package metrics

import (
	"fmt"
	"strings"
	"time"

	"github.com/VictoriaMetrics/metrics"
	"github.com/mr-karan/logchef/pkg/models"
)

// Metrics system with meaningful labels for monitoring LogChef usage

// RecordHTTPRequest records HTTP request metrics with optional user context
func RecordHTTPRequest(method, endpoint string, statusCode int, duration time.Duration, responseSize int64, user *models.User) {
	// HTTP request metrics
	labels := fmt.Sprintf(`logchef_http_requests_total{method="%s",endpoint="%s",status="%d"}`, method, endpoint, statusCode)
	metrics.GetOrCreateCounter(labels).Inc()
	
	// HTTP metrics with user context if available
	if user != nil {
		userLabels := fmt.Sprintf(`logchef_http_requests_by_user_total{method="%s",endpoint="%s",status="%d",user_email="%s",user_role="%s"}`, 
			method, endpoint, statusCode, user.Email, string(user.Role))
		metrics.GetOrCreateCounter(userLabels).Inc()
	}
	
	durationLabels := fmt.Sprintf(`logchef_http_request_duration_seconds{method="%s",endpoint="%s"}`, method, endpoint)
	metrics.GetOrCreateHistogram(durationLabels).Update(duration.Seconds())
	
	sizeLabels := fmt.Sprintf(`logchef_http_response_size_bytes{method="%s",endpoint="%s"}`, method, endpoint)
	metrics.GetOrCreateHistogram(sizeLabels).Update(float64(responseSize))
}

// RecordHTTPError records HTTP error metrics
func RecordHTTPError(method, endpoint, errorType string, statusCode int, user *models.User) {
	labels := fmt.Sprintf(`logchef_http_errors_total{method="%s",endpoint="%s",error_type="%s",status="%d"}`, method, endpoint, errorType, statusCode)
	metrics.GetOrCreateCounter(labels).Inc()
	
	if user != nil {
		userLabels := fmt.Sprintf(`logchef_http_errors_by_user_total{method="%s",endpoint="%s",error_type="%s",status="%d",user_email="%s",user_role="%s"}`, 
			method, endpoint, errorType, statusCode, user.Email, string(user.Role))
		metrics.GetOrCreateCounter(userLabels).Inc()
	}
}

// RecordQuery records query execution metrics with source context
func RecordQuery(source *models.Source, queryType string, success bool, duration time.Duration, rowsReturned int64, user *models.User) {
	result := "success"
	if !success {
		result = "failure"
	}
	
	// Query metrics with meaningful source labels
	sourceLabels := fmt.Sprintf(`logchef_query_total{source_id="%d",source_name="%s",database="%s",table="%s",query_type="%s",result="%s"}`, 
		source.ID, source.Name, source.Connection.Database, source.Connection.TableName, queryType, result)
	metrics.GetOrCreateCounter(sourceLabels).Inc()
	
	// User-specific query metrics if user context available
	if user != nil {
		userLabels := fmt.Sprintf(`logchef_query_by_user_total{source_name="%s",database="%s",table="%s",query_type="%s",result="%s",user_email="%s",user_role="%s"}`, 
			source.Name, source.Connection.Database, source.Connection.TableName, queryType, result, user.Email, string(user.Role))
		metrics.GetOrCreateCounter(userLabels).Inc()
	}
	
	// Duration histogram with source context
	durationLabels := fmt.Sprintf(`logchef_query_duration_seconds{source_name="%s",database="%s",table="%s"}`, 
		source.Name, source.Connection.Database, source.Connection.TableName)
	metrics.GetOrCreateHistogram(durationLabels).Update(duration.Seconds())
	
	if success && rowsReturned >= 0 {
		rowsLabels := fmt.Sprintf(`logchef_query_rows_returned{source_name="%s",database="%s",table="%s"}`, 
			source.Name, source.Connection.Database, source.Connection.TableName)
		metrics.GetOrCreateHistogram(rowsLabels).Update(float64(rowsReturned))
	}
}

// RecordQueryTimeout records query timeout metrics
func RecordQueryTimeout(source *models.Source, queryType string) {
	labels := fmt.Sprintf(`logchef_query_timeouts_total{source_id="%d",source_name="%s",database="%s",table="%s",query_type="%s"}`, 
		source.ID, source.Name, source.Connection.Database, source.Connection.TableName, queryType)
	metrics.GetOrCreateCounter(labels).Inc()
}

// RecordQueryError records query error metrics
func RecordQueryError(source *models.Source, errorType string) {
	labels := fmt.Sprintf(`logchef_query_errors_total{source_id="%d",source_name="%s",database="%s",table="%s",error_type="%s"}`, 
		source.ID, source.Name, source.Connection.Database, source.Connection.TableName, errorType)
	metrics.GetOrCreateCounter(labels).Inc()
}

// RecordHistogram records histogram generation metrics
func RecordHistogram(source *models.Source, success bool, duration time.Duration, user *models.User) {
	result := "success"
	if !success {
		result = "failure"
	}
	
	labels := fmt.Sprintf(`logchef_histogram_total{source_id="%d",source_name="%s",database="%s",table="%s",result="%s"}`, 
		source.ID, source.Name, source.Connection.Database, source.Connection.TableName, result)
	metrics.GetOrCreateCounter(labels).Inc()
	
	if user != nil {
		userLabels := fmt.Sprintf(`logchef_histogram_by_user_total{source_name="%s",database="%s",table="%s",result="%s",user_email="%s"}`, 
			source.Name, source.Connection.Database, source.Connection.TableName, result, user.Email)
		metrics.GetOrCreateCounter(userLabels).Inc()
	}
	
	durationLabels := fmt.Sprintf(`logchef_histogram_duration_seconds{source_name="%s",database="%s",table="%s"}`, 
		source.Name, source.Connection.Database, source.Connection.TableName)
	metrics.GetOrCreateHistogram(durationLabels).Update(duration.Seconds())
}

// RecordClickHouseConnectionStatus sets connection status for a source
func RecordClickHouseConnectionStatus(source *models.Source, healthy bool) {
	status := 0.0
	if healthy {
		status = 1.0
	}
	
	labels := fmt.Sprintf(`logchef_clickhouse_connection_status{source_id="%d",source_name="%s",database="%s",table="%s",host="%s"}`, 
		source.ID, source.Name, source.Connection.Database, source.Connection.TableName, source.Connection.Host)
	metrics.GetOrCreateGauge(labels, nil).Set(status)
}

// RecordClickHouseValidation records connection validation metrics
func RecordClickHouseValidation(source *models.Source, success bool) {
	result := "success"
	if !success {
		result = "failure"
	}
	
	labels := fmt.Sprintf(`logchef_clickhouse_connection_validation_total{source_id="%d",source_name="%s",database="%s",table="%s",host="%s",result="%s"}`, 
		source.ID, source.Name, source.Connection.Database, source.Connection.TableName, source.Connection.Host, result)
	metrics.GetOrCreateCounter(labels).Inc()
}

// RecordClickHouseReconnection records reconnection attempts
func RecordClickHouseReconnection(source *models.Source, success bool) {
	result := "success"
	if !success {
		result = "failure"
	}
	
	labels := fmt.Sprintf(`logchef_clickhouse_reconnections_total{source_id="%d",source_name="%s",database="%s",table="%s",host="%s",result="%s"}`, 
		source.ID, source.Name, source.Connection.Database, source.Connection.TableName, source.Connection.Host, result)
	metrics.GetOrCreateCounter(labels).Inc()
}

// RecordAuthAttempt records authentication attempt metrics
func RecordAuthAttempt(method string, success bool, user *models.User) {
	result := "success"
	if !success {
		result = "failure"
	}
	
	// Authentication metrics
	labels := fmt.Sprintf(`logchef_auth_attempts_total{method="%s",result="%s"}`, method, result)
	metrics.GetOrCreateCounter(labels).Inc()
	
	// Auth metrics with user context if available
	if user != nil {
		userLabels := fmt.Sprintf(`logchef_auth_attempts_by_user_total{method="%s",result="%s",user_email="%s",user_role="%s"}`, 
			method, result, user.Email, string(user.Role))
		metrics.GetOrCreateCounter(userLabels).Inc()
	}
}

// RecordSessionOperation records session operation metrics
func RecordSessionOperation(operation string, success bool, user *models.User) {
	result := "success"
	if !success {
		result = "failure"
	}
	
	labels := fmt.Sprintf(`logchef_session_operations_total{operation="%s",result="%s"}`, operation, result)
	metrics.GetOrCreateCounter(labels).Inc()
	
	if user != nil {
		userLabels := fmt.Sprintf(`logchef_session_operations_by_user_total{operation="%s",result="%s",user_email="%s",user_role="%s"}`, 
			operation, result, user.Email, string(user.Role))
		metrics.GetOrCreateCounter(userLabels).Inc()
	}
}

// RecordAPITokenOperation records API token operation metrics
func RecordAPITokenOperation(operation string, success bool, user *models.User, tokenName string) {
	result := "success"
	if !success {
		result = "failure"
	}
	
	labels := fmt.Sprintf(`logchef_api_token_operations_total{operation="%s",result="%s"}`, operation, result)
	metrics.GetOrCreateCounter(labels).Inc()
	
	if user != nil {
		userLabels := fmt.Sprintf(`logchef_api_token_operations_by_user_total{operation="%s",result="%s",user_email="%s",token_name="%s"}`, 
			operation, result, user.Email, tokenName)
		metrics.GetOrCreateCounter(userLabels).Inc()
	}
}

// RecordAuthorizationFailure records authorization failure metrics
func RecordAuthorizationFailure(endpoint string, user *models.User, reason string) {
	labels := fmt.Sprintf(`logchef_authorization_failures_total{endpoint="%s",reason="%s"}`, endpoint, reason)
	metrics.GetOrCreateCounter(labels).Inc()
	
	if user != nil {
		userLabels := fmt.Sprintf(`logchef_authorization_failures_by_user_total{endpoint="%s",reason="%s",user_email="%s",user_role="%s"}`, 
			endpoint, reason, user.Email, string(user.Role))
		metrics.GetOrCreateCounter(userLabels).Inc()
	}
}

// RecordTeamOperation records team operation metrics
func RecordTeamOperation(team *models.Team, operation string, success bool, user *models.User) {
	result := "success"
	if !success {
		result = "failure"
	}
	
	labels := fmt.Sprintf(`logchef_team_operations_total{team_id="%d",team_name="%s",operation="%s",result="%s"}`, 
		team.ID, team.Name, operation, result)
	metrics.GetOrCreateCounter(labels).Inc()
	
	if user != nil {
		userLabels := fmt.Sprintf(`logchef_team_operations_by_user_total{team_name="%s",operation="%s",result="%s",user_email="%s",user_role="%s"}`, 
			team.Name, operation, result, user.Email, string(user.Role))
		metrics.GetOrCreateCounter(userLabels).Inc()
	}
}

// RecordSourceOperation records source lifecycle operation metrics
func RecordSourceOperation(source *models.Source, operation string, success bool, user *models.User) {
	result := "success"
	if !success {
		result = "failure"
	}
	
	labels := fmt.Sprintf(`logchef_source_operations_total{source_id="%d",source_name="%s",database="%s",table="%s",operation="%s",result="%s"}`, 
		source.ID, source.Name, source.Connection.Database, source.Connection.TableName, operation, result)
	metrics.GetOrCreateCounter(labels).Inc()
	
	if user != nil {
		userLabels := fmt.Sprintf(`logchef_source_operations_by_user_total{source_name="%s",database="%s",table="%s",operation="%s",result="%s",user_email="%s",user_role="%s"}`, 
			source.Name, source.Connection.Database, source.Connection.TableName, operation, result, user.Email, string(user.Role))
		metrics.GetOrCreateCounter(userLabels).Inc()
	}
}

// RecordCollectionOperation records saved query collection operations
func RecordCollectionOperation(source *models.Source, team *models.Team, operation string, success bool, user *models.User, collectionName string) {
	result := "success"
	if !success {
		result = "failure"
	}
	
	labels := fmt.Sprintf(`logchef_collection_operations_total{source_name="%s",team_name="%s",operation="%s",result="%s",collection_name="%s"}`, 
		source.Name, team.Name, operation, result, collectionName)
	metrics.GetOrCreateCounter(labels).Inc()
	
	if user != nil {
		userLabels := fmt.Sprintf(`logchef_collection_operations_by_user_total{source_name="%s",team_name="%s",operation="%s",result="%s",user_email="%s"}`, 
			source.Name, team.Name, operation, result, user.Email)
		metrics.GetOrCreateCounter(userLabels).Inc()
	}
}

// RecordAIOperation records AI operation metrics
func RecordAIOperation(source *models.Source, operation string, success bool, duration time.Duration, user *models.User) {
	result := "success"
	if !success {
		result = "failure"
	}
	
	labels := fmt.Sprintf(`logchef_ai_operations_total{source_name="%s",database="%s",table="%s",operation="%s",result="%s"}`, 
		source.Name, source.Connection.Database, source.Connection.TableName, operation, result)
	metrics.GetOrCreateCounter(labels).Inc()
	
	if user != nil {
		userLabels := fmt.Sprintf(`logchef_ai_operations_by_user_total{source_name="%s",operation="%s",result="%s",user_email="%s"}`, 
			source.Name, operation, result, user.Email)
		metrics.GetOrCreateCounter(userLabels).Inc()
	}
	
	durationLabels := fmt.Sprintf(`logchef_ai_duration_seconds{source_name="%s",operation="%s"}`, source.Name, operation)
	metrics.GetOrCreateHistogram(durationLabels).Update(duration.Seconds())
}

// Gauge metrics for active counts
func UpdateActiveUsers(count int) {
	metrics.GetOrCreateGauge("logchef_active_users", nil).Set(float64(count))
}

func UpdateActiveSessions(count int) {
	metrics.GetOrCreateGauge("logchef_active_sessions", nil).Set(float64(count))
}

func IncrementActiveRequests() {
	metrics.GetOrCreateGauge("logchef_http_active_requests", nil).Inc()
}

func DecrementActiveRequests() {
	metrics.GetOrCreateGauge("logchef_http_active_requests", nil).Dec()
}

// Utility functions for query analysis

// DetermineQueryType analyzes the query string to determine its type
func DetermineQueryType(query string) string {
	query = strings.ToLower(strings.TrimSpace(query))
	
	if strings.HasPrefix(query, "select") {
		return "select"
	}
	if strings.HasPrefix(query, "show") {
		return "show"
	}
	if strings.HasPrefix(query, "describe") || strings.HasPrefix(query, "desc") {
		return "describe"
	}
	if strings.HasPrefix(query, "explain") {
		return "explain"
	}
	if strings.HasPrefix(query, "insert") {
		return "insert"
	}
	if strings.HasPrefix(query, "create") {
		return "create"
	}
	if strings.HasPrefix(query, "drop") {
		return "drop"
	}
	if strings.HasPrefix(query, "alter") {
		return "alter"
	}
	
	return "other"
}

// DetermineErrorType determines the type of error from an error object
func DetermineErrorType(err error) string {
	if err == nil {
		return ""
	}
	
	errStr := strings.ToLower(err.Error())
	
	if strings.Contains(errStr, "timeout") || strings.Contains(errStr, "deadline") {
		return "timeout"
	}
	if strings.Contains(errStr, "connection") || strings.Contains(errStr, "network") {
		return "connection"
	}
	if strings.Contains(errStr, "syntax") || strings.Contains(errStr, "parse") {
		return "syntax"
	}
	if strings.Contains(errStr, "permission") || strings.Contains(errStr, "access") {
		return "permission"
	}
	if strings.Contains(errStr, "not found") || strings.Contains(errStr, "missing") {
		return "not_found"
	}
	
	return "other"
}

// contains checks if a string contains a substring
func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

// trimAndLower trims whitespace and converts to lowercase
func trimAndLower(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}