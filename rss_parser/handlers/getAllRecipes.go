package handlers

// GetAllRecipes 获取存储的数据
// 非常重要:!!!!!!
//
// 在操作 MongoDB 的增删改查时, 要注意,所有的结构体对象的字段必须要可导出, 否则将无法解析到数据.
func GetAllRecipes(c *gin.Context) {
	collection := mongoClient.Database(utils.MONGO_DATABASE).Collection(utils.MONGO_RSS_PARSER_COLLECTION)
	cur, err := collection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	defer cur.Close(ctx)
	
	entries := make([]*struct {
		Title     string `bson:"title"`
		Thumbnail string `bson:"thumbnail"`
		Link      string `bson:"link"`
	}, 0)
	
	for cur.Next(ctx) {
		var entry struct {
			Title     string `bson:"title"`
			Thumbnail string `bson:"thumbnail"`
			Link      string `bson:"link"`
		}
		err := cur.Decode(&entry)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		entries = append(entries, &entry)
	}
	c.JSON(http.StatusOK, entries)
	
}
