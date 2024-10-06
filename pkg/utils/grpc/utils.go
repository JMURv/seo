package utils

import (
	pb "github.com/JMURv/seo-svc/api/pb"
	md "github.com/JMURv/seo-svc/pkg/model"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func PagesToProto(req []*md.Page) []*pb.PageMsg {
	var res []*pb.PageMsg
	for _, v := range req {
		res = append(res, PageToProto(v))
	}
	return res
}

func PageToProto(req *md.Page) *pb.PageMsg {
	return &pb.PageMsg{
		Slug:      req.Slug,
		Title:     req.Title,
		Href:      req.Href,
		CreatedAt: timestamppb.New(req.CreatedAt),
		UpdatedAt: timestamppb.New(req.UpdatedAt),
	}
}

func ProtoToPage(req *pb.PageMsg) *md.Page {
	return &md.Page{
		Slug:      req.Slug,
		Title:     req.Title,
		Href:      req.Href,
		CreatedAt: req.CreatedAt.AsTime(),
		UpdatedAt: req.UpdatedAt.AsTime(),
	}
}

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
