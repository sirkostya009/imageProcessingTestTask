package web

import (
	"context"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"golang.org/x/image/draw"
	"image"
	"image/jpeg"
	"imageProcessingTestTask/db"
	"os"
	"sync"
)

var volumePath = os.Getenv("VOLUME_PATH")

func resizeAndWriteImage(src image.Image, size int, uid int32, tx *db.Queries, lock *sync.Mutex, ctx context.Context) error {
	if errors.Is(ctx.Err(), context.Canceled) {
		return nil
	}
	dst := image.NewRGBA(image.Rect(0, 0, size, size))
	draw.ApproxBiLinear.Scale(dst, dst.Rect, src, src.Bounds(), draw.Over, nil)
	path := fmt.Sprintf("%s/%d/%d.jpg", volumePath, uid, size)
	lock.Lock()
	err := tx.CreateImage(ctx, db.CreateImageParams{
		UserID: uid,
		Url:    path,
	})
	lock.Unlock()
	if err != nil {
		return fmt.Errorf("failed to persist image: %w", err)
	}
	out, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	err = jpeg.Encode(out, dst.SubImage(dst.Rect), nil)
	if err != nil {
		return fmt.Errorf("failed to write image: %w", err)
	}
	return nil
}

func (h *Handlers) UploadImage(c echo.Context) error {
	f, err := c.FormFile("image")
	if err != nil {
		return echo.ErrBadRequest
	}

	file, err := f.Open()
	if err != nil {
		return echo.ErrInternalServerError
	}

	img, _, err := image.Decode(file)
	_ = file.Close()
	if err != nil {
		return echo.ErrUnsupportedMediaType
	}

	ctx := c.Request().Context()
	tx, err := h.db.BeginTx(ctx)
	if err != nil {
		return echo.ErrInternalServerError
	}
	defer tx.Rollback(ctx)

	qtx := h.db.WithTx(tx)
	uid := c.Get("userId").(int32)
	if qtx.DeleteOldImages(ctx, uid) != nil {
		return echo.ErrInternalServerError
	}

	err = os.MkdirAll(fmt.Sprintf("%s/%d", volumePath, uid), os.ModePerm)
	if err != nil {
		return echo.ErrInternalServerError
	}

	var wg sync.WaitGroup
	var lock sync.Mutex
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for i := range 5 {
		go func() {
			defer wg.Done()
			for j := range 10 {
				if err := resizeAndWriteImage(img, 30+i*10+j, uid, qtx, &lock, ctx); err != nil {
					fmt.Println(err)
					cancel()
					return
				}
			}
		}()
		wg.Add(1)
	}

	if err = resizeAndWriteImage(img, 80, uid, qtx, &lock, ctx); err != nil {
		fmt.Println(err)
		cancel()
	}

	wg.Wait()

	if ctx.Err() != nil {
		return echo.NewHTTPError(500, "images were not saved")
	} else {
		err = tx.Commit(ctx)
		if err != nil {
			return echo.ErrInternalServerError
		}
		return c.NoContent(201)
	}
}

func (h *Handlers) GetImages(c echo.Context) error {
	images, err := h.db.GetImages(c.Request().Context(), c.Get("userId").(int32))
	if err != nil {
		return echo.ErrInternalServerError
	}

	return c.JSON(200, images)
}
