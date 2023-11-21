package model

import "time"

type Deployment struct {
	ID             int        `gorm:"column:id;primaryKey;autoIncrement"`
	UID            int        `gorm:"column:uid;not null;default:0"`
	SpaceID        string     `gorm:"column:space_id;not null;default:''"`
	SpaceName      string     `gorm:"column:space_name;not null;default:''"`
	CfgName        string     `gorm:"column:cfg_name;not null;default:''"`
	Paid           string     `gorm:"column:paid;not null;default:''"`
	Duration       int        `gorm:"column:duration;not null;default:0"`
	TxHash         string     `gorm:"column:tx_hash;not null;default:''"`
	ChainID        string     `gorm:"column:chain_id;not null;default:''"`
	Region         string     `gorm:"column:region;not null;default:''"`
	StartIn        int        `gorm:"column:start_in;not null;default:0"`
	JobID          string     `gorm:"column:job_id;not null;default:''"`
	ResultURL      string     `gorm:"column:result_url;not null;default:''"`
	ProviderID     string     `gorm:"column:provider_id;size:256;not null;default:''"`
	ProviderNodeID string     `gorm:"column:provider_node_id;size:256;not null;default:''"`
	Cost           string     `gorm:"column:cost;not null;default:''"`
	Status         int        `gorm:"column:status;type:tinyint;not null;default:0"`
	StatusMsg      string     `gorm:"column:status_msg;size:64;not null;default:''"`
	Msg            string     `gorm:"column:msg;size:1024;not null;default:''"`
	Source         int        `gorm:"column:source;type:tinyint;default:0"`
	ExpiredAt      *time.Time `gorm:"column:expired_at;type:timestamp"`
	CreatedAt      time.Time  `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt      time.Time  `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
}

func (Deployment) TableName() string {
	return "deployment"
}
