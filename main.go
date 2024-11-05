package main

import (
	"fmt"
	"sync"
	"time"
	"math/rand"
)

// Estructura para representar una pista de aterrizaje
type Runway struct {
	useTime   int
	variation int
	isOccupied bool
	mutex     sync.Mutex
}

// Estructura para representar una puerta de desembarque
type Gate struct {
	useTime   int
	variation int
	isOccupied bool
	mutex     sync.Mutex
}

// Estructura para representar un avión
type Plane struct {
	id        int
	hasLanded bool
}

// Estructura para representar la torre de control
type ControlTower struct {
	runways    []Runway
	gates      []Gate
	maxWaitTime int
}

// Función para inicializar una pista
func initializeRunway(runway *Runway, useTime int, variation int) {
	runway.useTime = useTime
	runway.variation = variation
	runway.isOccupied = false
}

// Función para inicializar una puerta
func initializeGate(gate *Gate, useTime int, variation int) {
	gate.useTime = useTime
	gate.variation = variation
	gate.isOccupied = false
}

// Función para solicitar aterrizaje en una pista
func (tower *ControlTower) requestLanding(plane *Plane) int {
	for i := 0; i < len(tower.runways); i++ {
		tower.runways[i].mutex.Lock()
		if !tower.runways[i].isOccupied {
			tower.runways[i].isOccupied = true
			plane.hasLanded = true
			tower.runways[i].mutex.Unlock()
			return i // Retorna el índice de la pista asignada
		}
		tower.runways[i].mutex.Unlock()
	}
	return -1 // Si no hay pistas disponibles
}

// Función para liberar una pista
func (tower *ControlTower) releaseRunway(runwayIndex int) {
	tower.runways[runwayIndex].mutex.Lock()
	tower.runways[runwayIndex].isOccupied = false
	tower.runways[runwayIndex].mutex.Unlock()
}

// Función para asignar una puerta de desembarque
func (tower *ControlTower) assignGate(plane *Plane) int {
	for i := 0; i < len(tower.gates); i++ {
		tower.gates[i].mutex.Lock()
		if !tower.gates[i].isOccupied {
			tower.gates[i].isOccupied = true
			tower.gates[i].mutex.Unlock()
			return i // Retorna el índice de la puerta asignada
		}
		tower.gates[i].mutex.Unlock()
	}
	return -1 // Si no hay puertas disponibles
}

// Función para liberar una puerta
func (tower *ControlTower) releaseGate(gateIndex int) {
	tower.gates[gateIndex].mutex.Lock()
	tower.gates[gateIndex].isOccupied = false
	tower.gates[gateIndex].mutex.Unlock()
}

// Función principal para simular el proceso de aterrizaje y desembarque
func main() {
	rand.Seed(time.Now().UnixNano())

	// Inicialización de la torre de control, pistas y puertas
	tower := ControlTower{
		runways:    make([]Runway, 3), // Ejemplo con 3 pistas
		gates:      make([]Gate, 2),   // Ejemplo con 2 puertas
		maxWaitTime: 10,
	}

	// Inicialización de las pistas y puertas con tiempos de uso y variación
	for i := range tower.runways {
		initializeRunway(&tower.runways[i], 10, 2)
	}
	for i := range tower.gates {
		initializeGate(&tower.gates[i], 5, 1)
	}

	// Simulación de llegada de aviones
	var wg sync.WaitGroup
	planes := []Plane{
		{id: 1}, {id: 2}, {id: 3},
	}

	for i := range planes {
		wg.Add(1)
		go func(plane *Plane) {
			defer wg.Done()

			// Solicita aterrizaje
			runwayIndex := tower.requestLanding(plane)
			if runwayIndex != -1 {
				fmt.Printf("Avión %d aterrizó en la pista %d\n", plane.id, runwayIndex)

				// Espera el tiempo de uso de la pista antes de liberarla
				time.Sleep(time.Duration(tower.runways[runwayIndex].useTime+rand.Intn(tower.runways[runwayIndex].variation)) * time.Second)
				tower.releaseRunway(runwayIndex)

				// Asigna una puerta de desembarque
				gateIndex := tower.assignGate(plane)
				if gateIndex != -1 {
					fmt.Printf("Avión %d desembarca en la puerta %d\n", plane.id, gateIndex)

					// Espera el tiempo de uso de la puerta antes de liberarla
					time.Sleep(time.Duration(tower.gates[gateIndex].useTime+rand.Intn(tower.gates[gateIndex].variation)) * time.Second)
					tower.releaseGate(gateIndex)
				} else {
					fmt.Printf("No hay puertas disponibles para el avión %d\n", plane.id)
				}
			} else {
				fmt.Printf("No hay pistas disponibles para el avión %d\n", plane.id)
			}
		}(&planes[i])
	}

	wg.Wait()
	fmt.Println("Simulación completada.")
}
