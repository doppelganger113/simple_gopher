package storage

type Order string

const (
	OrderDescending Order = "DESC"
	OrderAscending  Order = "ASC"
)

const (
	PaginationLimitDefault  = 20
	PaginationLimitMax      = 50
	PaginationOffsetDefault = 0
)

func PagingToLimitOffset(page, size uint) (limit, offset int) {
	limit = PaginationLimitDefault
	offset = PaginationOffsetDefault

	if size == 0 {
		limit = PaginationLimitDefault
	} else if size > PaginationLimitDefault {
		limit = PaginationLimitMax
	} else {
		limit = int(size)
	}

	if page > 1 {
		offset = (int(page) - 1) * limit
	}

	return
}

func (order Order) IsValid() bool {
	switch order {
	case OrderAscending:
		return true
	case OrderDescending:
		return true
	}

	return false
}

func ToOrderOr(value string, defaultOrder Order) Order {
	orderType := Order(value)
	if orderType.IsValid() {
		return orderType
	}

	return defaultOrder
}
