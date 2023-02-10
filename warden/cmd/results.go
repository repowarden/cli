package cmd

import "fmt"

// an enum for what an auditResult can be
type auditResultType int

const (
	RESULT_DEBUG auditResultType = iota
	RESULT_INFO
	RESULT_WARNING
	RESULT_ERROR
)

// auditResults are what a policy audit can return. These can be notes, debug information, warnings, and most importantly, errors.
type auditResult struct {
	repository *wardenRepo
	resultType auditResultType
	message    string
	values     []any
}

// Properly print out a result
func (this auditResult) String() string {
	return fmt.Sprintf(this.message, this.values...)
}

// the list of results returned from the audit
type auditResults []auditResult

// adds a new result
func (this *auditResults) add(repo *wardenRepo, resultType auditResultType, message string, values ...any) {
	*this = append(*this, auditResult{
		repo,
		resultType,
		message,
		values,
	})
}

// Returns a subset, just the one type
func (this *auditResults) ByType(resultType auditResultType) auditResults {

	var results auditResults

	for _, result := range *this {

		if result.resultType == resultType {
			results = append(results, result)
		}
	}

	return results
}

// combine two auditResults together
func (this *auditResults) merge(results auditResults) {
	*this = append(*this, results...)
}
