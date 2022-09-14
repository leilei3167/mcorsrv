package model

type Category struct {
	//实际开发过程中,尽量设置为 not null
	//https://zhuanlan.zhihu.com/p/73997266
	BaseModel
	Name string `gorm:"type:varchar(20);not null"` //最长20字符,不为空

	ParentCategoryID int32     //其上一级的分类的主键
	ParentCategory   *Category //其上一级的分类,指向自身

	//这些类型是使用int32还是int? 应该尽量采用int32,和proto生成的文件保持一致,减少强制转换
	Level int32 `gorm:"type:int;not null;default:1;comment:'默认为1级分类'"`
	IsTab bool  `gorm:"default:false;not null"` //是否侧边栏显示
}

type Brands struct {
	BaseModel
	Name string `gorm:"type:varchar(20);not null"`             //最长20字符,不为空
	Logo string `gorm:"type:varchar(200);default:'';not null"` //图片是url 要长一些
	//品牌应该和分类要有关联,否则无法通过商品分类进行过滤,但不能是一对一的关系,因为一个品牌会有多种类型的产品
	//这就是典型的多对多的关系
}

// GoodsCategoryBrand 是一个辅助表,将分类和品牌以联合索引的形式关联起来
type GoodsCategoryBrand struct {
	BaseModel
	//应该建立唯一索引
	CategoryID int32 `gorm:"type:int;index:idx_category_brand,unique"`
	Category   Category

	BrandsID int32 `gorm:"type:int;index:idx_category_brand,unique"` //和CategoryID创建联合唯一索引
	Brands   Brands
}

// TableName 自定义表名
func (GoodsCategoryBrand) TableName() string {
	return "goodscategorybrand"
}

// Banner 轮播图,需要一张图片,以及对应的跳转至商品详细页的url,以及一个顺序 index
type Banner struct {
	BaseModel
	Image string `gorm:"type:varchar(200);not null"`
	Url   string `gorm:"type:varchar(20);not null"`
	Index int32  `gorm:"type:int;default:1;not null"`
}

type Goods struct {
	BaseModel

	CategoryID int32 `gorm:"type:int;not null"`
	Category   Category
	BrandsID   int32 `gorm:"type:int;not null"`
	Brands     Brands

	OnSale   bool `gorm:"default:false;not null"` //是否上架
	ShipFree bool `gorm:"default:false;not null"` //包邮?
	IsNew    bool `gorm:"default:false;not null"`
	IsHot    bool `gorm:"default:false;not null"`

	Name        string  `gorm:"type:varchar(50);not null"`   //商品名
	GooddSn     string  `gorm:"type:varchar(50);not null"`   //商品的产品编号,店家用
	ClickNum    int32   `gorm:"type:int;default:0;not null"` //点击量
	SoldNum     int32   `gorm:"type:int;default:0;not null"` //销量
	FavNum      int32   `gorm:"type:int;default:0;not null"` //收藏数
	MarketPrice float32 `gorm:"not null"`
	ShopPrice   float32 `gorm:"not null"`
	GoodsBrief  string  `gorm:"type:varchar(100);not null"` //简介

	//图片都是url形式,无法以切片表示,需要引入辅助表或者自定义类型
	//商品图片
	Images GormList `gorm:"type:varchar(1000);not null"`
	//详细图片
	DescImages GormList `gorm:"type:varchar(1000);not null"`
	//封面图
	GoodsFrontImage GormList `gorm:"type:varchar(1000);not null"`
}

/*// GoodsImages 添加一张外键表是可以的,但是图片量大了之后,join的性能会降低;字符串在数据库中搜索比较慢
type GoodsImages struct {
	GoodsId int32
	Image   string
}*/
