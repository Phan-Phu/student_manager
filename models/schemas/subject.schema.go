package schemas

type SubjectRequest struct {
	Name string `json:"name"`
}

func (*SubjectRequest) Validate() {

}

type SubjectResponse struct {
	ID         string `json:"subject_id" bson:"_id"`
	Name       string `json:"name" bson:"name"`
	CreateDate string `json:"create_date" bson:"create_date"`
}
