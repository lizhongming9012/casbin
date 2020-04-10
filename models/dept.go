package models

type Dept struct {
	ID       uint   `gorm:"primary_key"`
	Parentid uint   `json:"parentid"`
	Name     string `json:"name"`
}

func AddDept(data interface{}) error {
	if err := db.Table("dept").Create(data).Error; err != nil {
		return err
	}
	return nil
}

func UpdateDept(d *Dept) error {
	if err := db.Table("dept").Where("id=?", d.ID).Updates(d).Error; err != nil {
		return err
	}
	return nil
}

func IsParent(id uint) bool {
	var d Dept
	if err := db.Where("parentid=? ", id).First(&d).Error; err != nil {
		return false
	}
	return true
}

func DeleteDept(id uint) error {
	if err := db.Where("id=?", id).Delete(Dept{}).Error; err != nil {
		return err
	}
	return nil
}

type DeptTree struct {
	ID       uint        `json:"id"`
	Name     string      `json:"name"`
	Parentid uint        `json:"parentid"`
	Children []*DeptTree `json:"children"`
}

//获取部门树
func GetDeptTree(depId uint) ([]DeptTree, error) {
	var dept Dept
	err := db.Where("id=?", depId).First(&dept).Error
	if err != nil {
		return nil, err
	}
	perms := make([]DeptTree, 0)
	child := DeptTree{
		ID:       dept.ID,
		Name:     dept.Name,
		Parentid: dept.Parentid,
		Children: []*DeptTree{},
	}
	err = getDeptTreeNode(depId, &child)
	if err != nil {
		return nil, err
	}
	perms = append(perms, child)
	return perms, nil
}

//递归获取子节点
func getDeptTreeNode(parentId uint, tree *DeptTree) error {
	var perms []*Dept
	//根据父结点Id查询数据表，获取相应的子结点信息
	err := db.Where("parentid=?", parentId).Find(&perms).Error
	if err != nil {
		return err
	}
	for i := 0; i < len(perms); i++ {
		child := DeptTree{
			ID:       perms[i].ID,
			Name:     perms[i].Name,
			Parentid: perms[i].Parentid,
			Children: []*DeptTree{},
		}
		tree.Children = append(tree.Children, &child)
		err = getDeptTreeNode(perms[i].ID, &child)
	}
	return err
}
