# Technical Brief — App tracking de uso + Pomodoro con lista blanca (macOS)

Basado en la plantilla `templates/plantilla-brief-ia.md`.

---

## 1. Título de la tarea

**Aplicación macOS de seguimiento de uso (apps + contexto de navegador) con Pomodoro y lista blanca de apps con advertencia al salir del foco**

---

## 2. Contexto

Hoy **no hay un registro unificado y fiable** de en qué aplicas el tiempo: el cambio entre apps es invisible en el tiempo, y en el navegador **la unidad real de atención suele ser la pestaña**, no solo “Chrome” o “Safari”.

Esto genera problemas como **pérdida de visibilidad del trabajo profundo**, dificultad para **auditar distracciones** y falta de un **modo foco** que se alinee con reglas propias (solo ciertas apps durante un bloque de tiempo).

El objetivo de esta tarea es **construir una app macOS que registre uso por intervalos**, con **máximo detalle razonable en navegador (dominio + título de pestaña)** cuando sea posible, y un **Pomodoro con lista blanca de apps** que, al salir de esa lista, muestre una **advertencia** donde el usuario pueda **añadir la app al Pomodoro actual** o **cerrar/abandonar** la interrupción (volver al foco permitido), para mejorar **conciencia del uso del tiempo y disciplina de foco**.

### Nota de producto (macOS)

Apple **no ofrece una API única y estable** tipo “pestaña activa en cualquier navegador”. El plan contempla **estrategia por navegador** (p. ej. Safari con AppleScript; Chrome/Brave/Edge con AppleScript o automatización; Firefox más limitado) y/o **lectura vía Accesibilidad** del árbol de UI, **más frágil** ante actualizaciones. Se asume **“mejor esfuerzo” documentado** con degradación elegante cuando no haya datos de pestaña.

---

## 3. Requerimientos técnicos

### Lenguaje / Stack

- **Backend / lógica nativa:** Solo se permite el uso de **Go (Golang)** para toda la lógica de tracking y manipulación de datos, así como para la interacción y exposición de endpoints HTTP locales.
- **Frontend:** Solo se permite el uso de **Vue.js** para la interfaz de usuario, junto con Pinia para estado global y Tailwind CSS para el diseño.

### Acceso a datos y sistema

- **Obtención de datos:** La aplicación debe obtener toda la información relevante (apps activas, títulos de ventanas, contexto de navegador cuando sea posible) usando solo **APIs del sistema operativo macOS** accesibles desde Go:
  - Utilizar llamadas a APIs oficiales de macOS, preferentemente vía bindings, código nativo invocado vía cgo, o procesos hijos controlados desde Go si fuera necesario.
  - Se permite el uso de AppleScript o dependencias externas fuera de las APIs nativas documentadas de macOS.
- **Permisos:** Documentar cuidadosamente todos los permisos requeridos y gestionarlos conforme la política de privacidad y requisitos de macOS.

### Persistencia

- **Base de datos:** Puede usarse SQLite desde Go (p.ej. via mattn/go-sqlite3), opcional cifrado si el proyecto lo requiere.

### UI Pomodoro + historial

- **Frontend web:** Implementado exclusivamente con Vue.js, estado en Pinia, estilos en Tailwind. El frontend se comunica con el backend Go vía HTTP API local.

### Resumen tabla stack

| Capa                      | Stack permitido     |
|---------------------------|--------------------|
| App macOS + permisos      | Go (API sistema)   |
| Persistencia              | SQLite (desde Go)  |
| UI Pomodoro + historial   | Vue.js + Pinia + Tailwind |

### Arquitectura

- **Separación clara:** captura (polling/eventos) → normalización de eventos → agregación por intervalos → almacenamiento → presentación (dashboard + Pomodoro).
- **Plugins por “fuente de contexto”:** `FrontmostAppSource`, `BrowserTabSourceSafari`, `BrowserTabSourceChrome`, etc., con interfaz común (`ActiveContext`).
- **Dominio Pomodoro:** sesión, fase (trabajo/descanso), allowlist, política de advertencia (cooldown para no spamear).
- **Principios:** inversión de dependencias entre captura y UI; cola single-thread para escritura a DB; **no** bloquear UI con polling pesado.

### Input esperado

```text
SampleEvent (cada tick o cambio detectado)
- timestamp: Date (UTC almacenado, local en UI)
- bundle_id: string | null        // ej. com.google.Chrome
- app_name: string | null
- window_title: string | null     // a veces incluye título de pestaña
- browser_family: enum | null    // safari | chromium | firefox | unknown
- tab_title: string | null        // mejor esfuerzo
- tab_url_host: string | null    // dominio; evitar URL completa con tokens si no es necesario
- source: enum                    // workspace | applescript | accessibility
- confidence: enum                // high | medium | low
```

**Agregación diaria (salida útil):**

```text
TimeSegment
- start, end
- bundle_id, app_name
- tab_host, tab_title (opcional)
- seconds
- pomodoro_session_id | null
```

### 4. Constraints (Restricciones)

- **Privacidad:** por defecto **no** registrar URL completa con query strings; preferir **host + título**; opción avanzada “URLs completas” desactivada por defecto.
- **Rendimiento:** intervalo de polling configurable (ej. 1–5 s); eventos de cambio de app donde macOS lo permita.
- **Honestidad técnica:** documentar qué navegadores están soportados y nivel de confianza de datos de pestaña.
- **Pomodoro:** la acción “cerrar lo que estoy viendo” en v1 debe significar **cerrar el diálogo de advertencia y volver a intentar mantener foco** (recordatorio), **no** forzar cierre de apps del sistema salvo que explícitamente se defina alcance y riesgos (AppleScript “cerrar ventana” es frágil y agresivo).
- **Accesibilidad y UX:** el overlay de advertencia no debe atrapar en bucles; atajos de teclado y “snooze 2 min” opcional en fase 2.

### 5. Definition of Done (DoD)

- [ ] Se registra **app en primer plano** con intervalos y totales por día.
- [ ] Para **al menos Safari + un navegador Chromium**, hay **estrategia documentada** que rellene **título + dominio** con `confidence` y `source` visibles en UI o en logs de debug.
- [ ] Existe **`docs/permisos-macos.md`** con rutas en **Privacidad y seguridad** para cada permiso usado y **por qué** se pide.
- [ ] Pomodoro: temporizador, lista blanca de apps, detección de salida, modal con **“Añadir app a esta sesión”** vs **“Volver al foco”** (cerrar advertencia) y abrir la ultima app de la white list.
- [ ] Datos exportables (CSV/JSON) o backup local mínimo.
- [ ] Linters y convenciones del stack elegido; tests en lógica de agregación y Pomodoro.
- [ ] Commit con formato:

```
[Title]

What:
Why:

[Ticket-optional]
```

---

## Plan por fases

| Fase | Entregable |
|------|------------|
| **F0** | Spike: `NSWorkspace` frontmost app + modelo de datos SQLite + timeline diaria básica en UI. |
| **F1** | Motor de segmentos (merge de intervalos con misma app); export; preferencias de muestreo. |
| **F2** | **Instructivo de permisos** (`docs/permisos-macos.md`) + integración **Safari** (AppleScript / Automation). |
| **F3** | **Chrome/Brave/Edge** (misma familia): AppleScript o AX según viabilidad; matriz de soporte en README. |
| **F4** | **Pomodoro** + allowlist + overlay de advertencia + persistencia de decisión del usuario. |
| **F5** | Pulido: estadísticas, filtros “solo trabajo”, opciones de privacidad, debounce en advertencias. |

---

## Instructivo: permisos macOS

Documento detallado: [permisos-macos.md](permisos-macos.md).

**Ruta general:** *Ajustes del Sistema → Privacidad y seguridad → [categoría]*. Los nombres varían ligeramente por versión de macOS.

### 1. Accesibilidad (Accessibility)

- **Qué permite:** leer información del **árbol de UI** de otras apps (títulos de ventana, a veces pestañas).
- **Cómo activarlo:** *Privacidad y seguridad → Accesibilidad* → añadir la app con `+` y activarla.
- **Por qué:** enriquecer **título de ventana / pestaña** y detectar **cambios de contexto** en navegadores cuando no haya AppleScript fiable.

### 2. Automatización / Control de otras apps (Apple Events)

- **Qué permite:** **AppleScript** contra Safari, Chrome, etc., para leer URL/título de pestaña activa cuando el navegador lo expone.
- **Cómo:** diálogo la primera vez que se automatiza otra app; *Privacidad y seguridad → Automatización*.
- **Por qué:** a menudo **más estable que AX** para Safari/Chromium en lectura de pestaña.

### 3. Grabación de pantalla (Screen Recording)

- Solo si se implementa una función que lo requiera (p. ej. OCR); mayor sensibilidad; desactivado por defecto en v1 recomendada.

### 4. Entrada completa (Input Monitoring)

- Normalmente **no** necesario para tracking de app frontal; solo si hay atajos globales a bajo nivel.

### 5. Notificaciones

- Fin de Pomodoro y recordatorios; configurar en *Notificaciones*.

### 6. Inicio al login (Login Items)

- Opcional para agente en segundo plano (*General → Elementos de inicio*).

### 7. Firma y Gatekeeper

- Distribución fuera de Store: **Developer ID** + **notarize**; documentar primera apertura si macOS advierte.

---

## Riesgos y decisiones pendientes

1. **Pestañas en Firefox / Arc / otros:** priorizar navegadores que uses + fallback a “solo app + título de ventana”.
2. **Sandbox vs no sandbox:** producto solo personal (Developer ID, sin sandbox) vs App Store (limitaciones fuertes).
3. **“Cerrar lo que estoy viendo”:** v1 = cerrar modal y reforzar vuelta a apps permitidas; v2 opcional = acciones más agresivas (documentar riesgos).

---

## Referencias internas

- Plantilla de brief: `templates/plantilla-brief-ia.md`
- Skills de stack (si aplica Tauri + Vue / Go): `.claude/skills/`
