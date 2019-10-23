package postage

const STICK_BOUNDARY = 10000000
const DISPLAY_INDEX_ORDER_ASC = 1
const DISPLAY_INDEX_ORDER_DESC = 2

//type itemPos struct {
//	Id int
//	Table string
//	DisplayIndex int
//	OriginalDisplayIndex int
//}
//
//var _VALID_ACTIONS = map[string]bool {
//	"up": true,
//	"down": true,
//	"top": true,
//	"bottom": true,
//	"stick_top": true,
//	"stick_bottom": true,
//	"unstick": true,
//}
//
//type UpdateDisplayIndexService struct {
//	eel.ServiceBase
//	order int
//}
//
//func NewUpdateDisplayIndexService(ctx context.Context, order int) *UpdateDisplayIndexService {
//	service := new(UpdateDisplayIndexService)
//	service.Ctx = ctx
//	service.order = order
//	return service
//}
//
//func getExtraFilterCondition(extraFilters eel.Map) string {
//	if extraFilters == nil || len(extraFilters) == 0 {
//		return "1 = 1"
//	}
//
//	buf := make([]string, 0)
//	for k, v := range extraFilters {
//		k = strings.Replace(k, "?", "%v", -1)
//		item := fmt.Sprintf(k, v)
//		buf = append(buf, item)
//	}
//
//	return strings.Join(buf, " and ")
//}
//
////doUp 向上调整
//func (this *UpdateDisplayIndexService) doUp(item *itemPos, extraFilters eel.Map) error {
//	o := eel.GetOrmFromContext(this.Ctx)
//	table := item.Table
//	targetItem := itemPos{}
//
//	//获得与item进行位置交换的targetItem
//	sql := fmt.Sprintf("select id, display_index, original_display_index from %s where %s and is_deleted = false and display_index < ? order by display_index desc limit 1", table, getExtraFilterCondition(extraFilters))
//	err := o.Raw(sql, item.DisplayIndex).QueryRow(&targetItem)
//	if err != nil {
//		eel.Logger.Error(err)
//		return err
//	}
//
//	if item.Id == targetItem.Id {
//		return nil
//	}
//
//	//如果targetItem是置顶或置底的，直接退出
//	if targetItem.DisplayIndex < -STICK_BOUNDARY {
//		return nil
//	}
//
//	//交换item与target的display_index
//	sql = fmt.Sprintf("update %s set display_index = ? where id = ?", table)
//	//调整targetItem的display_index
//	_, err = o.Raw(sql, item.DisplayIndex, targetItem.Id).Exec()
//	if err != nil {
//		eel.Logger.Error(err)
//		return err
//	}
//	//调整item的display_index
//	_, err = o.Raw(sql, targetItem.DisplayIndex, item.Id).Exec()
//	if err != nil {
//		eel.Logger.Error(err)
//		return err
//	}
//
//	return nil
//}
//
////doDown 向下调整
//func (this *UpdateDisplayIndexService) doDown(item *itemPos, extraFilters eel.Map) error {
//	o := eel.GetOrmFromContext(this.Ctx)
//	table := item.Table
//	targetItem := itemPos{}
//
//	//获得与item进行位置交换的targetItem
//	sql := fmt.Sprintf("select id, display_index, original_display_index from %s where %s and is_deleted = false and display_index > ? order by display_index limit 1", table, getExtraFilterCondition(extraFilters))
//	err := o.Raw(sql, item.DisplayIndex).QueryRow(&targetItem)
//	if err != nil {
//		eel.Logger.Error(err)
//		return err
//	}
//
//	if item.Id == targetItem.Id {
//		return nil
//	}
//
//	//如果targetItem是置顶或置底的，直接退出
//	if targetItem.DisplayIndex > STICK_BOUNDARY {
//		return nil
//	}
//
//	//交换item与target的display_index
//	sql = fmt.Sprintf("update %s set display_index = ? where id = ?", table)
//	//调整targetItem的display_index
//	_, err = o.Raw(sql, item.DisplayIndex, targetItem.Id).Exec()
//	if err != nil {
//		eel.Logger.Error(err)
//		return err
//	}
//	//调整item的display_index
//	_, err = o.Raw(sql, targetItem.DisplayIndex, item.Id).Exec()
//	if err != nil {
//		eel.Logger.Error(err)
//		return err
//	}
//
//	return nil
//}
//
////setTop 将item对应的记录移动到最顶端
//func (this *UpdateDisplayIndexService) setTop(item *itemPos, extraFilters eel.Map) error {
//	o := eel.GetOrmFromContext(this.Ctx)
//	table := item.Table
//	targetItem := itemPos{}
//
//	//获得排在第一位的targetItem
//	sql := fmt.Sprintf("select id, display_index, original_display_index from %s where %s and is_deleted = false and display_index > ? order by display_index limit 1", table, getExtraFilterCondition(extraFilters))
//	err := o.Raw(sql, -STICK_BOUNDARY).QueryRow(&targetItem)
//	if err != nil {
//		eel.Logger.Error(err)
//		return err
//	}
//
//	if item.Id == targetItem.Id {
//		return nil
//	}
//
//	//将item的排序设置在targetItem之前
//	sql = fmt.Sprintf("update %s set display_index = ? where id = ?", table)
//	_, err = o.Raw(sql, targetItem.DisplayIndex-1, item.Id).Exec()
//	if err != nil {
//		eel.Logger.Error(err)
//		return err
//	}
//
//	return nil
//}
//
////setBottom 将item对应的记录移动到最底端
//func (this *UpdateDisplayIndexService) setBottom(item *itemPos, extraFilters eel.Map) error {
//	o := eel.GetOrmFromContext(this.Ctx)
//	table := item.Table
//	targetItem := itemPos{}
//
//	//获得排在第一位的targetItem
//	sql := fmt.Sprintf("select id, display_index, original_display_index from %s where %s and is_deleted = false and display_index < ? order by display_index desc limit 1", table, getExtraFilterCondition(extraFilters))
//	err := o.Raw(sql, STICK_BOUNDARY).QueryRow(&targetItem)
//	if err != nil {
//		eel.Logger.Error(err)
//		return err
//	}
//
//	if item.Id == targetItem.Id {
//		return nil
//	}
//
//	//将item的排序设置在targetItem之前
//	sql = fmt.Sprintf("update %s set display_index = ? where id = ?", table)
//	_, err = o.Raw(sql, targetItem.DisplayIndex+1, item.Id).Exec()
//	if err != nil {
//		eel.Logger.Error(err)
//		return err
//	}
//
//	return nil
//}
//
////stickTop 将item对应的记录置顶
//func (this *UpdateDisplayIndexService) stickTop(item *itemPos, extraFilters eel.Map) error {
//	o := eel.GetOrmFromContext(this.Ctx)
//	table := item.Table
//
//	//将item的排序设置在targetItem之前
//	sql := fmt.Sprintf("update %s set display_index = ?, original_display_index = ? where id = ?", table)
//	_, err := o.Raw(sql, item.DisplayIndex-2*STICK_BOUNDARY, item.DisplayIndex, item.Id).Exec()
//	if err != nil {
//		eel.Logger.Error(err)
//		return err
//	}
//
//	return nil
//}
//
////stickTop 将item对应的记录置底
//func (this *UpdateDisplayIndexService) stickBottom(item *itemPos, extraFilters eel.Map) error {
//	o := eel.GetOrmFromContext(this.Ctx)
//	table := item.Table
//
//	//将item的排序设置在targetItem之前
//	sql := fmt.Sprintf("update %s set display_index = ?, original_display_index = ? where id = ?", table)
//	_, err := o.Raw(sql, item.DisplayIndex+2*STICK_BOUNDARY, item.DisplayIndex, item.Id).Exec()
//	if err != nil {
//		eel.Logger.Error(err)
//		return err
//	}
//
//	return nil
//}
//
////unstick 取消置顶或置底
//func (this *UpdateDisplayIndexService) unstick(item *itemPos, extraFilters eel.Map) error {
//	o := eel.GetOrmFromContext(this.Ctx)
//	table := item.Table
//
//	if !(item.DisplayIndex > STICK_BOUNDARY || item.DisplayIndex < -STICK_BOUNDARY) {
//		//不在置顶或置底的范围内，直接返回
//		return nil
//	}
//
//	sql := fmt.Sprintf("update %s set display_index = ?, original_display_index = 0 where id = ?", table)
//	_, err := o.Raw(sql, item.OriginalDisplayIndex, item.Id).Exec()
//	if err != nil {
//		eel.Logger.Error(err)
//		return err
//	}
//
//	return nil
//}
//
//func (this *UpdateDisplayIndexService) isAscOrder() bool {
//	return this.order == DISPLAY_INDEX_ORDER_ASC
//}
//
////Encode 对单个实体对象进行编码
//func (this *UpdateDisplayIndexService) Update(item *itemPos, action string) error {
//	return this.UpdateWithExtraFilter(item, action, nil)
//}
//
////Encode 对单个实体对象进行编码
//func (this *UpdateDisplayIndexService) UpdateWithExtraFilter(item *itemPos, action string, extraFilters eel.Map) error {
//	if _, ok := _VALID_ACTIONS[action]; !ok {
//		return errors.New(fmt.Sprintf("无效的排序操作:%s", action))
//	}
//
//	switch action {
//	case "up":
//		if this.isAscOrder() {
//			return this.doUp(item, extraFilters)
//		} else {
//			return this.doDown(item, extraFilters)
//		}
//	case "down":
//		if this.isAscOrder() {
//			return this.doDown(item, extraFilters)
//		} else {
//			return this.doUp(item, extraFilters)
//		}
//	case "top":
//		if this.isAscOrder() {
//			return this.setTop(item, extraFilters)
//		} else {
//			return this.setBottom(item, extraFilters)
//		}
//	case "bottom":
//		if this.isAscOrder() {
//			return this.setBottom(item, extraFilters)
//		} else {
//			return this.setTop(item, extraFilters)
//		}
//	case "stick_top":
//		if this.isAscOrder() {
//			return this.stickTop(item, extraFilters)
//		} else {
//			return this.stickBottom(item, extraFilters)
//		}
//	case "stick_bottom":
//		if this.isAscOrder() {
//			return this.stickBottom(item, extraFilters)
//		} else {
//			return this.stickTop(item, extraFilters)
//		}
//	case "unstick":
//		return this.unstick(item, extraFilters)
//	}
//
//	return nil
//}


func init() {
}
