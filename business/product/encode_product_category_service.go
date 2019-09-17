package product

import (
	"context"
	"fmt"
	"github.com/gingerxman/eel"
)

type EncodeProductCategoryService struct {
	eel.ServiceBase
}

func NewEncodeProductCategoryService(ctx context.Context) *EncodeProductCategoryService {
	service := new(EncodeProductCategoryService)
	service.Ctx = ctx
	return service
}

//Encode 对单个实体对象进行编码
func (this *EncodeProductCategoryService) Encode(productCategory *ProductCategory) *RProductCategory {
	if productCategory == nil {
		return nil
	}

	return &RProductCategory{
		Id: productCategory.Id,
		Name: productCategory.Name,
		NodeType: productCategory.NodeType,
		ProductCount: productCategory.ProductCount,
		DisplayIndex: productCategory.DisplayIndex,
		IsEnabled: productCategory.IsEnabled,
		CreatedAt: productCategory.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}

//EncodeMany 对实体对象进行批量编码
func (this *EncodeProductCategoryService) EncodeMany(product_categories []*ProductCategory) []*RProductCategory {
	rDatas := make([]*RProductCategory, 0)
	for _, productCategory := range product_categories {
		rDatas = append(rDatas, this.Encode(productCategory))
	}
	
	return rDatas
}

//EncodeTree 编码为tree
func (this *EncodeProductCategoryService) EncodeTree(product_categories []*ProductCategory) *RProductCategoryTreeNode {
	id2node := make(map[int]*RProductCategoryTreeNode)
	//add root node
	id2node[0] = &RProductCategoryTreeNode{
		Id: 0,
		Name: "root",
		Level: 0,
		SubCategories: make([]*RProductCategoryTreeNode, 0),
	}
	for _, productCategory := range product_categories {
		if node, ok := id2node[productCategory.Id]; !ok {
			level := -1
			fatherId := productCategory.FatherId
			var fatherNode *RProductCategoryTreeNode
			var ok2 = false
			if fatherNode, ok2 = id2node[fatherId]; ok2 {
				level = fatherNode.Level + 1
			} else {
				panic(eel.NewBusinessError("product_category_tree:invalid_father_id", fmt.Sprintf("no father_id(%d) in id2node", fatherId)))
			}
			
			//没有记录，创建node
			node = &RProductCategoryTreeNode{
				Id: productCategory.Id,
				Name: productCategory.Name,
				FatherId: fatherId,
				Level: level,
				SubCategories: make([]*RProductCategoryTreeNode, 0),
			}
			id2node[node.Id] = node
			fatherNode.SubCategories = append(fatherNode.SubCategories, node)
		}
	}
	
	return id2node[0]
}

func init() {
}
