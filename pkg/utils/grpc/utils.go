package utils

import (
	pb "github.com/JMURv/par-pro-seo/api/pb"
	md "github.com/JMURv/par-pro-seo/pkg/model"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ModelToProto(req *md.SEO) *pb.SEOMsg {
	return &pb.SEOMsg{
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

func ProtoToModel(req *pb.SEOMsg) *md.SEO {
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
