package package_calculator

import (
	"fmt"
	"sort"
	"database/sql"
)

type TopPackageCount struct {
	IndirectCount int
	Deps          []string
}

func CalculatePackages(rows *sql.Rows) error {
	relationships := make([][2]string, 0)

	for rows.Next() {
		var fromPackage, toPackage string
		if err := rows.Scan(&fromPackage, &toPackage); err != nil {
			return err
		}
		relationships = append(relationships, [2]string{fromPackage, toPackage})
	}

	topPackageMap := make(map[string]map[string]struct{})

	for _, relationship := range relationships {
		from := relationship[0]
		to := relationship[1]

		if _, exists := topPackageMap[to]; !exists {
			topPackageMap[to] = make(map[string]struct{})
		}
		topPackageMap[to][from] = struct{}{}
	}

	topPackageCounts := make([]TopPackageCount, 0)

	var dfs func(pkg string, visited map[string]struct{}, packages map[string]map[string]struct{}) []string
	dfs = func(pkg string, visited map[string]struct{}, packages map[string]map[string]struct{}) []string {
		if _, seen := visited[pkg]; seen {
			return nil
		}
		visited[pkg] = struct{}{}
		var allDeps []string
		if deps, exists := packages[pkg]; exists {
			for dep := range deps {
				allDeps = append(allDeps, dep)
				allDeps = append(allDeps, dfs(dep, visited, packages)...)
			}
		}
		return allDeps
	}

	for pkg := range topPackageMap {
		visited := make(map[string]struct{})
		deps := dfs(pkg, visited, topPackageMap)
		indirectCount := len(deps)
		topPackageCounts = append(topPackageCounts, TopPackageCount{IndirectCount: indirectCount, Deps: deps})
	}

	sort.Slice(topPackageCounts, func(i, j int) bool {
		return topPackageCounts[i].IndirectCount > topPackageCounts[j].IndirectCount
	})

	uniqueFromPackages := make(map[string]struct{})

	for _, fromPackages := range topPackageMap {
		for fromPackage := range fromPackages {
			uniqueFromPackages[fromPackage] = struct{}{}
		}
	}

	totalUniqueFromPackages := len(uniqueFromPackages)
	fmt.Println("Total unique frompackages:", totalUniqueFromPackages)
	threshold := int(float64(totalUniqueFromPackages) * 0.7)

	currentCount := 0
	requiredPackages := 0

	for _, pkgCount := range topPackageCounts {
		deps := pkgCount.Deps

		for _, dep := range deps {
			if _, exists := uniqueFromPackages[dep]; exists {
				currentCount++
				delete(uniqueFromPackages, dep)
			}
			if currentCount >= threshold {
				requiredPackages++
				break
			}
		}
		if currentCount >= threshold {
			break
		}
		requiredPackages++
	}

	fmt.Println("Number of remaining packages:", len(uniqueFromPackages))
	fmt.Printf("It requires the top %d packages to cover 70%% of unique frompackages.\n", requiredPackages)
	return nil
}