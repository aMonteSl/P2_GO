package main

import (
	"testing"
	"sync"
)

// Test básico: simula el sistema con valores nominales
func TestBasicSimulation(t *testing.T) {
	numAirplanes := 10
	numRunways := 3
	numGates := 5
	towerConfig := Config{Time: 100, StdTime: 20, Buffer: numRunways}
	runwayConfig := Config{Time: 200, StdTime: 50, Buffer: numRunways}
	gateConfig := Config{Time: 300, StdTime: 100, Buffer: numGates}

	success := runSimulation(numAirplanes, numRunways, numGates, towerConfig, runwayConfig, gateConfig)
	if !success {
		t.Errorf("Simulación básica fallida")
	}
}

// Test: doblar capacidad de espera en la torre de control
func TestDoubleControlTowerBuffer(t *testing.T) {
	numAirplanes := 10
	numRunways := 3
	numGates := 5
	towerConfig := Config{Time: 100, StdTime: 20, Buffer: 2 * numRunways}
	runwayConfig := Config{Time: 200, StdTime: 50, Buffer: numRunways}
	gateConfig := Config{Time: 300, StdTime: 100, Buffer: numGates}

	success := runSimulation(numAirplanes, numRunways, numGates, towerConfig, runwayConfig, gateConfig)
	if !success {
		t.Errorf("Simulación con buffer duplicado en torre fallida")
	}
}

// Test: incrementar tiempos de operación en un 25%
func TestIncreaseOperationTime(t *testing.T) {
	numAirplanes := 10
	numRunways := 3
	numGates := 5
	towerConfig := Config{Time: 125, StdTime: 25, Buffer: numRunways}
	runwayConfig := Config{Time: 250, StdTime: 62, Buffer: numRunways}
	gateConfig := Config{Time: 375, StdTime: 125, Buffer: numGates}

	success := runSimulation(numAirplanes, numRunways, numGates, towerConfig, runwayConfig, gateConfig)
	if !success {
		t.Errorf("Simulación con incremento de tiempos fallida")
	}
}

// Test: multiplicar el número de pistas por 5
func TestIncreaseRunways(t *testing.T) {
	numAirplanes := 10
	numRunways := 15
	numGates := 5
	towerConfig := Config{Time: 100, StdTime: 20, Buffer: numRunways}
	runwayConfig := Config{Time: 200, StdTime: 50, Buffer: numRunways}
	gateConfig := Config{Time: 300, StdTime: 100, Buffer: numGates}

	success := runSimulation(numAirplanes, numRunways, numGates, towerConfig, runwayConfig, gateConfig)
	if !success {
		t.Errorf("Simulación con aumento de pistas fallida")
	}
}

// Test: multiplicar pistas por 5 y aumentar su tiempo 5 veces
func TestIncreaseRunwaysWithTime(t *testing.T) {
	numAirplanes := 10
	numRunways := 15
	numGates := 5
	towerConfig := Config{Time: 100, StdTime: 20, Buffer: numRunways}
	runwayConfig := Config{Time: 1000, StdTime: 250, Buffer: numRunways}
	gateConfig := Config{Time: 300, StdTime: 100, Buffer: numGates}

	success := runSimulation(numAirplanes, numRunways, numGates, towerConfig, runwayConfig, gateConfig)
	if !success {
		t.Errorf("Simulación con aumento de pistas y tiempos fallida")
	}
}

// Función auxiliar para ejecutar la simulación
func runSimulation(numAirplanes, numRunways, numGates int, towerConfig, runwayConfig, gateConfig Config) bool {
	defer func() {
		if r := recover(); r != nil {
			// Recuperar de un pánico en caso de error crítico
		}
	}()

	// Inicialización
	runwayChan := make(chan *Runway, towerConfig.Buffer)
	wg := &sync.WaitGroup{}
	gateChan := make(chan *Gate, gateConfig.Buffer)

	for i := 1; i <= numGates; i++ {
		gateChan <- &Gate{id: i, config: gateConfig}
	}
	for i := 1; i <= numRunways; i++ {
		runwayChan <- &Runway{id: i, gateChan: gateChan, config: runwayConfig}
	}

	controlTower := &ControlTower{runwayChan: runwayChan, config: towerConfig, wg: wg}

	// Simulación de aviones
	for i := 1; i <= numAirplanes; i++ {
		wg.Add(1)
		airplane := &Airplane{id: i}
		go controlTower.handleAirplane(airplane)
	}

	wg.Wait()

	// Simulación terminada
	return true
}
