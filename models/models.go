package models

type Book struct {
	ISBN            int    `json:"isbn"`      
	LibID           int    `json:"libId"`
	Title           string `json:"title"`
	Author          string `json:"author"`
	Publisher       string `json:"publisher"`
	Version         string `json:"version"`
	TotalCopies     int    `json:"totalCopies"`
	AvailableCopies int    `json:"availableCopies"`
}

type User struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	Email         string `json:"email"`
	ContactNumber int    `json:"contactNumber"`
	Role          string `json:"role"`
	LibId         int    `json:"libId"`
}
type LibraryAdmin struct {
	ID   int    `json:"Id"`
	Name string `json:"name"`
}
type IssueRegistery struct {
	IssueID            int    `json:"issueId"`
	ISBN               int    `json:"isbn"`            //fk book
	ReaderID           int    `json:"readerId"`        //fk users
	IssueApproverID    int    `json:"issueApproverId"` //fk admin
	IssueStatus        string `json:"issueStatus"`
	IssueDate          string `json:"issueDate"`
	ExpectedReturnDate string `json:"expectedReturnDate"`
	ReturnDate         string `json:"returnDate"`
	ReturnApproverID   int    `json:"returnApproverId"` //fk admin

}
type RequestEvents struct {
	ReqID        int       `json:"reqId"`
	BookID       int       `json:"bookId"`   //fk book
	ReaderID     int       `json:"readerId"` //fk user
	RequestDate  string `json:"requestDate"`
	ApprovalDate string `json:"approvalDate"`
	ApproverID   int       `json:"approverId"` //fk admin
	RequestType  string    `json:"requestType"`
}

