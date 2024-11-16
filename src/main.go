package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Configuración general
type Config struct {
	Time    int // Tiempo base de operación
	StdTime int // Variación en el tiempo
	Buffer  int // Capacidad máxima de espera (buffer)
}

// Avión
type Airplane struct {
	id int
}

// Torre de control
type ControlTower struct {
	runwayChan chan *Runway
	config     Config
	wg         *sync.WaitGroup
}

// Pista
type Runway struct {
	id      int
	gateChan chan *Gate
	config   Config
}

// Puerta de desembarque
type Gate struct {
	id      int
	config  Config
}

// Simulación del tiempo con variación
func simulateTime(baseTime, stdTime int) {
	duration := time.Duration(baseTime+rand.Intn(2*stdTime+1)-stdTime) * time.Millisecond
	time.Sleep(duration)
}

// Torre de Control: gestión de aviones hacia pistas
func (ct *ControlTower) handleAirplane(airplane *Airplane) {
	defer ct.wg.Done()

	fmt.Printf("Avión %d: Solicita pista...\n", airplane.id)
	runway := <-ct.runwayChan
	fmt.Printf("Avión %d: Asignada pista %d.\n", airplane.id, runway.id)
	simulateTime(ct.config.Time, ct.config.StdTime)

	// Avión pasa a la pista para aterrizar
	runway.handleAirplane(airplane)

	// Liberar pista
	ct.runwayChan <- runway
	fmt.Printf("Avión %d: Liberó pista %d.\n", airplane.id, runway.id)
}

// Pista: gestión de aviones hacia puertas
func (r *Runway) handleAirplane(airplane *Airplane) {
	fmt.Printf("Avión %d: Aterrizando en pista %d...\n", airplane.id, r.id)
	simulateTime(r.config.Time, r.config.StdTime)
	fmt.Printf("Avión %d: Aterrizó en pista %d. Solicita puerta...\n", airplane.id, r.id)

	// Solicita puerta
	gate := <-r.gateChan
	fmt.Printf("Avión %d: Asignada puerta %d.\n", airplane.id, gate.id)

	// Puerta desembarca pasajeros
	gate.handleAirplane(airplane)

	// Liberar puerta
	r.gateChan <- gate
	fmt.Printf("Avión %d: Liberó puerta %d.\n", airplane.id, gate.id)
}

// Puerta: desembarque
func (g *Gate) handleAirplane(airplane *Airplane) {
	fmt.Printf("Avión %d: Desembarcando en puerta %d...\n", airplane.id, g.id)
	simulateTime(g.config.Time, g.config.StdTime)
	fmt.Printf("Avión %d: Pasajeros desembarcados en puerta %d.\n", airplane.id, g.id)

	// Simular despegue
	fmt.Printf("Avión %d: Despegando tras completar desembarque en puerta %d...\n", airplane.id, g.id)
	simulateTime(g.config.Time, g.config.StdTime)
	fmt.Printf("Avión %d: Despegó exitosamente.\n", airplane.id)
}

func main() {
	rand.Seed(time.Now().UnixNano())

	// Configuración
	numAirplanes := 10   // Número de aviones
	numRunways := 3      // Número de pistas
	numGates := 5        // Número de puertas

	towerConfig := Config{Time: 100, StdTime: 20, Buffer: numRunways} // Torre de control
	runwayConfig := Config{Time: 200, StdTime: 50, Buffer: numRunways} // Pistas
	gateConfig := Config{Time: 300, StdTime: 100, Buffer: numGates}    // Puertas

	// Inicialización de la torre de control
	runwayChan := make(chan *Runway, towerConfig.Buffer)
	wg := &sync.WaitGroup{}

	// Inicialización de pistas y puertas
	gateChan := make(chan *Gate, gateConfig.Buffer)
	for i := 1; i <= numGates; i++ {
		gateChan <- &Gate{id: i, config: gateConfig}
	}
	for i := 1; i <= numRunways; i++ {
		runwayChan <- &Runway{id: i, gateChan: gateChan, config: runwayConfig}
	}

	// Torre de control
	controlTower := &ControlTower{runwayChan: runwayChan, config: towerConfig, wg: wg}

	// Simulación de aviones
	for i := 1; i <= numAirplanes; i++ {
		wg.Add(1)
		airplane := &Airplane{id: i}
		go controlTower.handleAirplane(airplane)
	}

	wg.Wait()
	fmt.Println("Simulación completada.")
}
