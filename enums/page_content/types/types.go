package types

// Page Content Types
const (
	Text   = "text"
	Table  = "table"
	List   = "list"
	Link   = "link"
	Img    = "img" // img stands for image
	Doc    = "doc" // doc stands for document
	Script = "script"

	Other     = "other"
	OtherFile = "other_file" // OtherFile represents unknown file type
)

// All returns slice enums
func All() []string {
	return []string{
		Text,
		Table,
		List,
		Link,
		Img,
		Doc,
		Script,
		Other,
		OtherFile,
	}
}
