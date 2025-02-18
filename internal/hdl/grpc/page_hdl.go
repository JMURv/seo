package grpc

import (
	"context"
	"errors"
	pb "github.com/JMURv/seo/api/grpc/v1/gen"
	ctrl "github.com/JMURv/seo/internal/ctrl"
	hdl "github.com/JMURv/seo/internal/hdl"
	"github.com/JMURv/seo/internal/hdl/validation"
	utils "github.com/JMURv/seo/internal/models/mapper"
	metrics "github.com/JMURv/seo/internal/observability/metrics/prometheus"
	ot "github.com/opentracing/opentracing-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

func (h *Handler) ListPages(ctx context.Context, req *pb.EmptySEO) (*pb.ListPageRes, error) {
	const op = "page.ListPages.hdl"
	s, c := time.Now(), codes.OK
	span, ctx := ot.StartSpanFromContext(ctx, op)
	defer func() {
		span.Finish()
		metrics.ObserveRequest(time.Since(s), int(c), op)
	}()

	if req == nil {
		c = codes.InvalidArgument
		return nil, status.Errorf(c, hdl.ErrDecodeRequest.Error())
	}

	res, err := h.ctrl.ListPages(ctx)
	if err != nil {
		span.SetTag("error", true)
		c = codes.Internal
		return nil, status.Errorf(c, hdl.ErrInternal.Error())
	}
	return &pb.ListPageRes{
		Pages: utils.PagesToProto(res),
	}, nil
}

func (h *Handler) GetPage(ctx context.Context, req *pb.SlugSEO) (*pb.PageMsg, error) {
	const op = "page.GetPage.hdl"
	s, c := time.Now(), codes.OK
	span, ctx := ot.StartSpanFromContext(ctx, op)
	defer func() {
		span.Finish()
		metrics.ObserveRequest(time.Since(s), int(c), op)
	}()

	if req == nil || req.Slug == "" {
		c = codes.InvalidArgument
		return nil, status.Errorf(c, hdl.ErrDecodeRequest.Error())
	}

	res, err := h.ctrl.GetPage(ctx, req.Slug)
	if err != nil && errors.Is(err, ctrl.ErrNotFound) {
		c = codes.NotFound
		return nil, status.Errorf(c, err.Error())
	} else if err != nil {
		span.SetTag("error", true)
		c = codes.Internal
		return nil, status.Errorf(c, hdl.ErrInternal.Error())
	}
	return utils.PageToProto(res), nil
}

func (h *Handler) CreatePage(ctx context.Context, req *pb.PageMsg) (*pb.SlugSEO, error) {
	const op = "page.CreatePage.hdl"
	s, c := time.Now(), codes.OK
	span, ctx := ot.StartSpanFromContext(ctx, op)
	defer func() {
		span.Finish()
		metrics.ObserveRequest(time.Since(s), int(c), op)
	}()

	if req == nil {
		c = codes.InvalidArgument
		return nil, status.Errorf(c, hdl.ErrDecodeRequest.Error())
	}

	obj := utils.ProtoToPage(req)
	if err := validation.ValidatePage(obj); err != nil {
		c = codes.InvalidArgument
		return nil, status.Errorf(c, err.Error())
	}

	res, err := h.ctrl.CreatePage(ctx, obj)
	if err != nil && errors.Is(err, ctrl.ErrNotFound) {
		c = codes.NotFound
		return nil, status.Errorf(c, err.Error())
	} else if err != nil {
		span.SetTag("error", true)
		c = codes.Internal
		return nil, status.Errorf(c, hdl.ErrInternal.Error())
	}
	return &pb.SlugSEO{
		Slug: res,
	}, nil
}

func (h *Handler) UpdatePage(ctx context.Context, req *pb.PageWithSlugMsg) (*pb.EmptySEO, error) {
	const op = "page.UpdatePage.hdl"
	s, c := time.Now(), codes.OK
	span, ctx := ot.StartSpanFromContext(ctx, op)
	defer func() {
		span.Finish()
		metrics.ObserveRequest(time.Since(s), int(c), op)
	}()

	if req == nil || req.Slug == "" || req.Page == nil {
		c = codes.InvalidArgument
		return nil, status.Errorf(c, hdl.ErrDecodeRequest.Error())
	}

	obj := utils.ProtoToPage(req.Page)
	if err := validation.ValidatePage(obj); err != nil {
		c = codes.InvalidArgument
		return nil, status.Errorf(c, err.Error())
	}

	err := h.ctrl.UpdatePage(ctx, req.Slug, obj)
	if err != nil && errors.Is(err, ctrl.ErrNotFound) {
		c = codes.NotFound
		return nil, status.Errorf(c, err.Error())
	} else if err != nil {
		span.SetTag("error", true)
		c = codes.Internal
		return nil, status.Errorf(c, hdl.ErrInternal.Error())
	}
	return &pb.EmptySEO{}, nil
}

func (h *Handler) DeletePage(ctx context.Context, req *pb.SlugSEO) (*pb.EmptySEO, error) {
	const op = "page.DeletePage.hdl"
	s, c := time.Now(), codes.OK
	span, ctx := ot.StartSpanFromContext(ctx, op)
	defer func() {
		span.Finish()
		metrics.ObserveRequest(time.Since(s), int(c), op)
	}()

	if req == nil || req.Slug == "" {
		c = codes.InvalidArgument
		return nil, status.Errorf(c, hdl.ErrDecodeRequest.Error())
	}

	err := h.ctrl.DeletePage(ctx, req.Slug)
	if err != nil && errors.Is(err, ctrl.ErrNotFound) {
		c = codes.NotFound
		return nil, status.Errorf(c, err.Error())
	} else if err != nil {
		span.SetTag("error", true)
		c = codes.Internal
		return nil, status.Errorf(c, hdl.ErrInternal.Error())
	}
	return &pb.EmptySEO{}, nil
}
