package constants

type RepositoryType string

const (
	RepositoryBeginTx       RepositoryType = "beginTx"
	RepositoryRollback      RepositoryType = "rollback"
	RepositoryCommit        RepositoryType = "commit"
	RepositoryPing          RepositoryType = "ping"
	RepositoryFind          RepositoryType = "find"
	RepositoryFindAll       RepositoryType = "findAll"
	RepositoryFindOne       RepositoryType = "findOne"
	RepositoryFindByID      RepositoryType = "findByID"
	RepositorySave          RepositoryType = "save"
	RepositoryUpdate        RepositoryType = "update"
	RepositoryDelete        RepositoryType = "delete"
	RepositoryCount         RepositoryType = "count"
	RepositoryCountDistinct RepositoryType = "countDistinct"
	RepositoryVerify        RepositoryType = "verify"
	RepositoryRaw           RepositoryType = "raw"
	RepositorySaveOrUpdate  RepositoryType = "saveOrUpdate"
	RepositoryBulkSave      RepositoryType = "bulkSave"
	RepositoryBulkUpdate    RepositoryType = "bulkUpdate"
	RepositoryBulkDelete    RepositoryType = "bulkDelete"
)
