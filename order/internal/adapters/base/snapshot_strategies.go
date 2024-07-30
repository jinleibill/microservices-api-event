package base

var DefaultSnapshotStrategies = NewMaxChangesSnapshotStrategy(10)

type SnapshotStrategy interface {
	ShouldSnapshot(root *AggregateRoot) bool
}

type maxChangesSnapshotStrategy struct {
	maxChanges int
}

func NewMaxChangesSnapshotStrategy(maxChanges int) SnapshotStrategy {
	return &maxChangesSnapshotStrategy{maxChanges: maxChanges}
}

func (m *maxChangesSnapshotStrategy) ShouldSnapshot(root *AggregateRoot) bool {
	return root.PendingVersion() >= m.maxChanges && len(root.Events()) >= m.maxChanges ||
		root.PendingVersion()%m.maxChanges < len(root.Events()) ||
		root.PendingVersion()%m.maxChanges == 0
}
