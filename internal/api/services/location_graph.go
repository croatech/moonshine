package services

import (
	"sync"

	"moonshine/internal/repository"

	"github.com/google/uuid"
)

type LocationGraph struct {
	adjacency map[uuid.UUID][]uuid.UUID
	slugToID  map[string]uuid.UUID
	idToSlug  map[uuid.UUID]string
	mu        sync.RWMutex
}

func NewLocationGraph(locationRepo *repository.LocationRepository) (*LocationGraph, error) {
	graph := &LocationGraph{
		adjacency: make(map[uuid.UUID][]uuid.UUID),
		slugToID:  make(map[string]uuid.UUID),
		idToSlug:  make(map[uuid.UUID]string),
	}

	locations, err := locationRepo.FindAll()
	if err != nil {
		return nil, err
	}

	for _, loc := range locations {
		graph.slugToID[loc.Slug] = loc.ID
		graph.idToSlug[loc.ID] = loc.Slug
	}

	connections, err := locationRepo.FindAllConnections()
	if err != nil {
		return nil, err
	}

	for _, conn := range connections {
		graph.adjacency[conn.LocationID] = append(graph.adjacency[conn.LocationID], conn.NearLocationID)
		graph.adjacency[conn.NearLocationID] = append(graph.adjacency[conn.NearLocationID], conn.LocationID)
	}

	return graph, nil
}

func (g *LocationGraph) FindShortestPath(fromSlug, toSlug string) ([]string, error) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	fromID, ok := g.slugToID[fromSlug]
	if !ok {
		return nil, repository.ErrLocationNotFound
	}

	toID, ok := g.slugToID[toSlug]
	if !ok {
		return nil, repository.ErrLocationNotFound
	}

	if fromID == toID {
		return []string{}, nil
	}

	queue := []uuid.UUID{fromID}
	visited := make(map[uuid.UUID]bool)
	parent := make(map[uuid.UUID]uuid.UUID)
	visited[fromID] = true

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if current == toID {
			path := g.reconstructPath(parent, fromID, toID)
			return path, nil
		}

		for _, neighbor := range g.adjacency[current] {
			if !visited[neighbor] {
				visited[neighbor] = true
				parent[neighbor] = current
				queue = append(queue, neighbor)
			}
		}
	}

	return nil, ErrLocationNotConnected
}

func (g *LocationGraph) reconstructPath(parent map[uuid.UUID]uuid.UUID, fromID, toID uuid.UUID) []string {
	path := []uuid.UUID{}
	current := toID

	for current != fromID {
		path = append([]uuid.UUID{current}, path...)
		current = parent[current]
	}

	result := make([]string, len(path))
	for i, id := range path {
		result[i] = g.idToSlug[id]
	}

	return result
}
