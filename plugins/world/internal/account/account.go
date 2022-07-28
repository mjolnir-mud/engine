package account

type Account struct {
	ID             string `bson:"_id,omitempty" json:"id"`
	Username       string `bson:"username" json:"username"`
	HashedPassword string `bson:"hashedPassword"`
	Email          string `bson:"email" json:"email"`
}
