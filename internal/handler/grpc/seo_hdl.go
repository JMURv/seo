package grpc

import (
	"context"
	"errors"
	pb "github.com/JMURv/seo-svc/api/pb"
	ctrl "github.com/JMURv/seo-svc/internal/controller"
	hdl "github.com/JMURv/seo-svc/internal/handler"
	metrics "github.com/JMURv/seo-svc/internal/metrics/prometheus"
	"github.com/JMURv/seo-svc/internal/validation"
	utils "github.com/JMURv/seo-svc/pkg/utils/grpc"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

func (h *Handler) GetSEO(ctx context.Context, req *pb.GetSEOReq) (*pb.SEOMsg, error) {
	const op = "seo.GetSEO.hdl"
	s, c := time.Now(), codes.OK
	span := opentracing.GlobalTracer().StartSpan(op)
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer func() {
		span.Finish()
		metrics.ObserveRequest(time.Since(s), int(c), op)
	}()

	if req == nil || req.Name == "" || req.Pk == "" {
		c = codes.InvalidArgument
		return nil, status.Errorf(c, hdl.ErrDecodeRequest.Error())
	}

	res, err := h.ctrl.GetSEO(ctx, req.Name, req.Pk)
	if err != nil && errors.Is(err, ctrl.ErrNotFound) {
		c = codes.NotFound
		return nil, status.Errorf(c, err.Error())
	} else if err != nil {
		span.SetTag("error", true)
		c = codes.Internal
		return nil, status.Errorf(c, hdl.ErrInternal.Error())
	}
	return utils.ModelToProto(res), nil
}

func (h *Handler) CreateSEO(ctx context.Context, req *pb.SEOMsg) (*pb.Uuid64SEO, error) {
	const op = "seo.CreateSEO.hdl"
	s, c := time.Now(), codes.OK
	span := opentracing.GlobalTracer().StartSpan(op)
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer func() {
		span.Finish()
		metrics.ObserveRequest(time.Since(s), int(c), op)
	}()

	obj := utils.ProtoToModel(req)
	if err := validation.ValidateSEO(obj); err != nil {
		c = codes.InvalidArgument
		return nil, status.Errorf(c, err.Error())
	}

	res, err := h.ctrl.CreateSEO(ctx, obj)
	if err != nil && errors.Is(err, ctrl.ErrNotFound) {
		c = codes.NotFound
		return nil, status.Errorf(c, err.Error())
	} else if err != nil {
		span.SetTag("error", true)
		c = codes.Internal
		return nil, status.Errorf(c, hdl.ErrInternal.Error())
	}
	return &pb.Uuid64SEO{
		Id: res,
	}, nil
}

func (h *Handler) UpdateSEO(ctx context.Context, req *pb.SEOMsg) (*pb.EmptySEO, error) {
	const op = "seo.UpdateSEO.hdl"
	s, c := time.Now(), codes.OK
	span := opentracing.GlobalTracer().StartSpan(op)
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer func() {
		span.Finish()
		metrics.ObserveRequest(time.Since(s), int(c), op)
	}()

	obj := utils.ProtoToModel(req)
	if err := validation.ValidateSEO(obj); err != nil {
		c = codes.InvalidArgument
		return nil, status.Errorf(c, err.Error())
	}

	err := h.ctrl.UpdateSEO(ctx, obj)
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

func (h *Handler) DeleteSEO(ctx context.Context, req *pb.GetSEOReq) (*pb.EmptySEO, error) {
	const op = "seo.DeleteSEO.hdl"
	s, c := time.Now(), codes.OK
	span := opentracing.GlobalTracer().StartSpan(op)
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer func() {
		span.Finish()
		metrics.ObserveRequest(time.Since(s), int(c), op)
	}()

	if req == nil || req.Name == "" || req.Pk == "" {
		c = codes.InvalidArgument
		return nil, status.Errorf(c, hdl.ErrDecodeRequest.Error())
	}

	err := h.ctrl.DeleteSEO(ctx, req.Name, req.Pk)
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
