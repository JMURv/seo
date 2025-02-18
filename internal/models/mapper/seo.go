package mapper

import (
	"github.com/JMURv/seo/api/grpc/v1/gen"
	md "github.com/JMURv/seo/internal/models"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ModelToProto(req *md.SEO) *gen.SEOMsg {
	return &gen.SEOMsg{
		Id:            req.ID,
		Title:         req.Title,
		Description:   req.Description,
		Keywords:      req.Keywords,
		OGTitle:       req.OGTitle,
		OGDescription: req.OGDescription,
		OGImage:       req.OGImage,
		ObjName:       req.OBJName,
		ObjPk:         req.OBJPK,
		CreatedAt:     timestamppb.New(req.CreatedAt),
		UpdatedAt:     timestamppb.New(req.UpdatedAt),
	}
}

func ProtoToModel(req *gen.SEOMsg) *md.SEO {
	return &md.SEO{
		ID:            req.Id,
		Title:         req.Title,
		Description:   req.Description,
		Keywords:      req.Keywords,
		OGTitle:       req.OGTitle,
		OGDescription: req.OGDescription,
		OGImage:       req.OGImage,
		OBJName:       req.ObjName,
		OBJPK:         req.ObjPk,
		CreatedAt:     req.CreatedAt.AsTime(),
		UpdatedAt:     req.UpdatedAt.AsTime(),
	}
}
