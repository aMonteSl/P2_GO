package main

import (
	"testing"
	"time"
)

// Test básico para configuración inicial de la torre de control
func baseSimulation(tower *ControlTower, planes []Plane, t *testing.T) {
	for i := range planes {
		// Solicita pista para aterrizar
		runwayIndex := tower.requestLanding(&planes[i])
		if runwayIndex == -1 {
			t.Errorf("Failed to assign runway to plane %d", planes[i].id)
			continue
		}
		time.Sleep(time.Millisecond * 10) // Simula tiempo en pista
		tower.releaseRunway(runwayIndex)

		// Asigna puerta para desembarcar
		gateIndex := tower.assignGate(&planes[i])
		if gateIndex == -1 {
			t.Errorf("Failed to assign gate to plane %d", planes[i].id)
			continue
		}
		time.Sleep(time.Millisecond * 10) // Simula tiempo en puerta
		tower.releaseGate(gateIndex)
	}
}

// Caso 1: La cantidad máxima de aviones esperando se duplica
func TestDoubleMaxWaitTime(t *testing.T) {
	tower := ControlTower{
		runways:     []Runway{{useTime: 10, variation: 2}, {useTime: 10, variation: 2}},
		gates:       []Gate{{useTime: 5, variation: 1}, {useTime: 5, variation: 1}},
		maxWaitTime: 10, // Duplicamos el tiempo de espera
	}

	planes := []Plane{{id: 1}, {id: 2}, {id: 3}, {id: 4}, {id: 5}}
	baseSimulation(&tower, planes, t)
}

// Caso 2: La variación en el tiempo de uso/utilización es del 25% por encima del nominal
func TestIncreaseVariation(t *testing.T) {
	tower := ControlTower{
		runways:     []Runway{{useTime: 10, variation: 3}, {useTime: 10, variation: 3}}, // Variación aproximada
		gates:       []Gate{{useTime: 5, variation: 1}, {useTime: 5, variation: 1}}, // Variación aproximada
		maxWaitTime: 5,
	}

	planes := []Plane{{id: 1}, {id: 2}, {id: 3}}
	baseSimulation(&tower, planes, t)
}

// Caso 3: La cantidad máxima de aviones esperando se duplica y la variación en el tiempo de uso es del 25% por encima del nominal
func TestDoubleWaitAndIncreaseVariation(t *testing.T) {
	tower := ControlTower{
		runways:     []Runway{{useTime: 10, variation: 3}, {useTime: 10, variation: 3}}, // Variación aproximada
		gates:       []Gate{{useTime: 5, variation: 1}, {useTime: 5, variation: 1}}, // Variación aproximada
		maxWaitTime: 10, // Duplicamos el tiempo de espera
	}

	planes := []Plane{{id: 1}, {id: 2}, {id: 3}, {id: 4}, {id: 5}}
	baseSimulation(&tower, planes, t)
}

// Caso 4: Las pistas se multiplican por 5
func TestIncreaseRunways(t *testing.T) {
	tower := ControlTower{
		runways:     make([]Runway, 10), // 5 veces más pistas
		gates:       []Gate{{useTime: 5, variation: 1}, {useTime: 5, variation: 1}},
		maxWaitTime: 5,
	}

	for i := range tower.runways {
		tower.runways[i] = Runway{useTime: 10, variation: 2}
	}

	planes := []Plane{{id: 1}, {id: 2}, {id: 3}, {id: 4}, {id: 5}}
	baseSimulation(&tower, planes, t)
}

// Caso 5: Las pistas se multiplican por 5, pero tardan 5 veces más del tiempo de uso cada una
func TestIncreaseRunwaysAndUsageTime(t *testing.T) {
	tower := ControlTower{
		runways:     make([]Runway, 10), // 5 veces más pistas
		gates:       []Gate{{useTime: 5, variation: 1}, {useTime: 5, variation: 1}},
		maxWaitTime: 5,
	}

	for i := range tower.runways {
		tower.runways[i] = Runway{useTime: 50, variation: 2} // Pistas con 5 veces más tiempo de uso
	}

	planes := []Plane{{id: 1}, {id: 2}, {id: 3}, {id: 4}, {id: 5}}
	baseSimulation(&tower, planes, t)
}
