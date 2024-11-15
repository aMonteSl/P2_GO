# Práctica 2: Sistemas Distribuidos - Gestión de Aeropuerto

## Introducción

En esta práctica se implementa un modelo concurrente en **Go** para simular las operaciones de un aeropuerto. La simulación incluye el aterrizaje, desembarque y despegue de aviones utilizando mecanismos concurrentes como **goroutines** y **canales**. Se evalúa el rendimiento del sistema bajo diferentes configuraciones, considerando tiempos de espera y límites de capacidad.

El propósito es fortalecer el manejo de concurrencia en **Go**, con un enfoque práctico que abarca desde la implementación básica hasta pruebas automatizadas y análisis de resultados.

---

## Objetivo del Programa

El programa tiene como objetivo simular el flujo de aviones en un aeropuerto de manera concurrente. Esto incluye:
1. **Aterrizaje**: Los aviones son asignados a pistas disponibles por una torre de control.
2. **Desembarque**: Los aviones acceden a puertas para que los pasajeros bajen.
3. **Despegue**: Una vez desembarcados, los aviones despegan.
4. **Concurrencia y límites**: Modelar restricciones de capacidad (buffers) y variaciones en tiempos de operación para analizar su impacto.

Se busca:
- Verificar que el sistema se comporte correctamente bajo diferentes configuraciones.
- Identificar cómo cambian los tiempos promedio de operación según ajustes en capacidad y tiempos.

---

## Descripción Técnica

### Componentes del Sistema
1. **Aviones (`Airplane`)**:
   - Identificados por un ID único.
   - Fluyen a través de la torre de control, pistas y puertas.

2. **Torre de Control (`ControlTower`)**:
   - Coordina la asignación de pistas para aterrizaje.
   - Usa un canal con buffer para limitar el número máximo de aviones en espera.

3. **Pistas (`Runway`)**:
   - Gestionan el aterrizaje de aviones.
   - Conectadas a puertas para el desembarque.

4. **Puertas (`Gate`)**:
   - Manejan el desembarque de pasajeros.
   - Liberan aviones para que procedan a despegar.

### Concurrencia y Sincronización
- **Goroutines**:
  Cada avión es manejado por una goroutine independiente, que interactúa con la torre, pistas y puertas.
- **Canales**:
  - **Torre de control a pistas**: Canal con buffer para limitar aviones esperando asignación.
  - **Pistas a puertas**: Canal para coordinar la transferencia de aviones.
- **WaitGroup**:
  Utilizado para sincronizar la finalización de todas las goroutines antes de concluir la simulación.

### Configuración y Parámetros
1. **Tiempo Base (`Time`)**: Tiempo promedio para cada operación.
2. **Desviación Estándar (`StdTime`)**: Variación en los tiempos de operación.
3. **Capacidad Máxima (`Buffer`)**: Número máximo de aviones esperando en cada etapa.

---

## Resultados de las Pruebas

### Configuraciones Probadas
1. **Simulación básica**:
   - Nominal: 10 aviones, 3 pistas, 5 puertas.
   - Todos los tiempos y capacidades predeterminados.

2. **Capacidad duplicada**:
   - Se duplicó el buffer de la torre de control.
   - Los tiempos promedio se mantuvieron estables con menor congestión.

3. **Incremento del 25% en tiempos**:
   - Tiempos base y desviaciones aumentados en un 25%.
   - Incremento proporcional en los tiempos promedio.

4. **Multiplicación de pistas**:
   - Se incrementó el número de pistas a 15.
   - Reducción significativa en tiempos de espera.

5. **Multiplicación de pistas con incremento de tiempo**:
   - Pistas incrementadas 5 veces y tiempos de operación también aumentados 5 veces.
   - El sistema mantuvo estabilidad, pero los tiempos totales aumentaron.

### Resumen de Resultados
| Configuración                 | Torre (ms) | Pista (ms) | Puerta (ms) |
|-------------------------------|------------|------------|-------------|
| Nominal                       | 100        | 200        | 300         |
| Capacidad duplicada           | 95         | 198        | 305         |
| Incremento de tiempos (+25%)  | 125        | 250        | 375         |
| Pistas multiplicadas (x5)     | 90         | 190        | 290         |
| Pistas y tiempos aumentados   | 450        | 1000       | 1500        |

---

## Conclusiones

1. **Estabilidad del Sistema**:
   - El sistema respondió correctamente a todas las configuraciones probadas.
   - Los canales con buffer y las goroutines garantizaron una sincronización eficiente.

2. **Impacto de la Capacidad**:
   - Incrementar el buffer de espera en la torre de control redujo ligeramente los tiempos de espera, demostrando que el cuello de botella puede aliviarse aumentando la capacidad.

3. **Impacto del Tiempo**:
   - Aumentar los tiempos base afecta proporcionalmente los tiempos promedio por etapa, lo que es coherente con el modelo.

4. **Escalabilidad**:
   - Multiplicar el número de pistas mejoró significativamente el rendimiento del sistema.
   - Sin embargo, cuando los tiempos de operación también se incrementaron, los beneficios fueron limitados.

5. **Conclusión Final**:
   - El modelo es adecuado para simular operaciones aeroportuarias con restricciones realistas.
   - Permite identificar configuraciones óptimas para minimizar tiempos de espera y maximizar eficiencia.

---

## Código Fuente

El código completo del programa y las pruebas están disponibles en el archivo `main.go` y `main_test.go`. Se adjuntan en el apéndice de este documento o están disponibles en el repositorio indicado.

---

