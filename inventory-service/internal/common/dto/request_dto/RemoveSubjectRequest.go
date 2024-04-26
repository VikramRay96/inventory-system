package request_dto

type RemoveSubjectRequest struct {
	SubjectId  string `json:"subject_id"`
	LessonName string `json:"lesson_name"`
}
