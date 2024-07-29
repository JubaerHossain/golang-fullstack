package persistence

import (
	"fmt"

	"github.com/JubaerHossain/cn-api/domain/categories/entity"
)

// getChildCategories retrieves the child categories for a given parent ID.
func (r *CategoryRepositoryImpl) getChildCategories(parentID uint64) ([]*entity.ResponseCategory, error) {
	query := "SELECT id, title, slug, `order`, status_id, parent_id FROM news_categories WHERE parent_id = ? ORDER BY `order` ASC"
	rows, err := r.app.MDB.Query(query, parentID)
	if err != nil {
		return nil, fmt.Errorf("database query error: %w", err)
	}
	defer rows.Close()

	childCategories := []*entity.ResponseCategory{}
	for rows.Next() {
		var category entity.ResponseCategory
		err := rows.Scan(&category.ID, &category.Title, &category.Slug, &category.Order, &category.StatusID, &category.ParentID)
		if err != nil {
			return nil, fmt.Errorf("rows scan error: %w", err)
		}
		childCategories = append(childCategories, &category)
		// childCategories, err := r.getChildCategories(category.ID)
		// if err != nil {
		// 	return nil, fmt.Errorf("get child categories error: %w", err)
		// }
		// category.ChildCategory = childCategories

	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}
	return childCategories, nil
}
