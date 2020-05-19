package binarysearchtree

type Comparator interface {
	/**
	 * @return 返回值等于0，代表e1和e2相等；返回值大于0，代表e1大于e2；返回值小于于0，代表e1小于e2
	 */
	Compare(e1 interface{}, e2 interface{}) int
}
