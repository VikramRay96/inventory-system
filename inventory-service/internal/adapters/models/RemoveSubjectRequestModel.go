package models

type RemoveSubjectRequestModel struct {
	SubjectId  string `json:"subject_id" bson:"subject_id"`
	LessonName string `json:"lesson_name" bson:"lesson_name"`
}
