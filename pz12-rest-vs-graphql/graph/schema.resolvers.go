package graph

import (
	"context"
	"fmt"

	"example.com/pz12-rest-vs-graphql/graph/model"
	"example.com/pz12-rest-vs-graphql/internal/store"
)

func (r *mutationResolver) CreateTask(ctx context.Context, input model.CreateTaskInput) (*model.Task, error) {

	task := &store.Task{
		ID:          fmt.Sprintf("t_%03d", len(store.Tasks)+1),
		Title:       input.Title,
		Description: input.Description,
		Done:        false,
	}

	store.Tasks = append(store.Tasks, task)

	return &model.Task{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		Done:        task.Done,
	}, nil
}

func (r *mutationResolver) UpdateTask(ctx context.Context, id string, input model.UpdateTaskInput) (*model.Task, error) {

	for _, t := range store.Tasks {

		if t.ID == id {

			if input.Done != nil {
				t.Done = *input.Done
			}

			return &model.Task{
				ID:          t.ID,
				Title:       t.Title,
				Description: t.Description,
				Done:        t.Done,
			}, nil
		}
	}

	return nil, fmt.Errorf("task not found")
}

func (r *queryResolver) Tasks(ctx context.Context) ([]*model.Task, error) {

	result := make([]*model.Task, 0, len(store.Tasks))

	for _, t := range store.Tasks {

		result = append(result, &model.Task{
			ID:          t.ID,
			Title:       t.Title,
			Description: t.Description,
			Done:        t.Done,
		})
	}

	return result, nil
}

func (r *queryResolver) Task(ctx context.Context, id string) (*model.Task, error) {

	for _, t := range store.Tasks {

		if t.ID == id {

			return &model.Task{
				ID:          t.ID,
				Title:       t.Title,
				Description: t.Description,
				Done:        t.Done,
			}, nil
		}
	}

	return nil, nil
}

func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}

func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }