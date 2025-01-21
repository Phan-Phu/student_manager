package services

import (
	"errors"
	"fmt"
	"student_manager/models"
	db "student_manager/models/db"
	"sync"
	"time"
)

func (repo *StudentRepoService) InsertManyStudentMultiThreads() error {
	start := time.Now()
	defer func() {
		elapsed := time.Since(start)
		fmt.Printf("InsertManyStudent took %s\n", elapsed)
	}()

	const totalStudents = 1000000
	const numThreads = 16
	batchSize := totalStudents / numThreads

	var wg sync.WaitGroup
	errChan := make(chan error, numThreads)
	threadTimes := make([]time.Duration, numThreads)

	// Worker function to process and insert students in batches
	worker := func(threadID, start, end int) {
		defer wg.Done()

		threadStart := time.Now()

		students := make([]db.Student, end-start)
		for i := start; i < end; i++ {
			// name := fmt.Sprintf("Name%d", i+1)
			// className := fmt.Sprintf("class %d", i+1)
			// student := db.NewStudent(name, className, "", 10, 0)
			// students[i-start] = *student
		}

		if err := repo.CreateMany(students); err != nil {
			errChan <- err
		}

		threadTimes[threadID] = time.Since(threadStart)
	}

	// Start worker threads
	for t := 0; t < numThreads; t++ {
		start := t * batchSize
		end := start + batchSize
		if t == numThreads-1 {
			end = totalStudents // Handle any remainder
		}

		wg.Add(1)
		go worker(t, start, end)
	}

	// Wait for all worker threads to complete
	wg.Wait()
	close(errChan)

	// Check for errors
	for err := range errChan {
		if err != nil {
			return errors.New(models.GetErrorMessage(models.ErrorCodeFailCreateStudent) + err.Error())
		}
	}

	// Print thread execution times
	for threadID, threadTime := range threadTimes {
		fmt.Printf("Thread %d took %s\n", threadID, threadTime)
	}

	return nil
}

func (repo *StudentRepoService) InsertManyStudentOneThreads() error {
	start := time.Now()

	students := []db.Student{}

	for i := 0; i < 1000; i++ {
		// name := fmt.Sprintf("Name%d", i+1)
		// className := fmt.Sprintf("class %d", i+1)
		// student := db.NewStudent(name, className, "", 10, 0)
		// students = append(students, *student)
	}

	err := repo.CreateMany(students)
	if err != nil {
		return fmt.Errorf(models.GetErrorMessage(models.ErrorCodeFailCreateStudent), err)
	}

	estTime := time.Since(start)
	fmt.Printf("InsertManyStudentOneThreads took %s\n", estTime)

	return nil
}
