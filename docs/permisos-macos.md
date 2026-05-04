# Permisos macOS — App de tracking y Pomodoro

Guía para usuarios y desarrolladores. Las rutas siguen **macOS Ventura y posteriores** (Ajustes del Sistema). En versiones antiguas puede decir “Preferencias del Sistema”.

---

## Dónde se gestionan

**Ajustes del Sistema → Privacidad y seguridad**

Algunos permisos aparecen solo **después** de que la app lo solicite por primera vez (diálogo del sistema).

---

## 1. Accesibilidad

**Ruta:** *Privacidad y seguridad → Accesibilidad*

**Qué permite:** que tu aplicación use las APIs de **Accesibilidad** para leer el árbol de UI de otras apps (por ejemplo título de ventana frontal, y en algunos casos información relacionada con pestañas del navegador).

**Cuándo lo necesita este proyecto:** enriquecer contexto cuando **AppleScript** no sea suficiente o falle; lectura de **título de ventana** que a veces incluye el título de la pestaña.

**Cómo añadir la app:**

1. Abre *Ajustes del Sistema → Privacidad y seguridad → Accesibilidad*.
2. Pulsa el botón **+** (o activa el interruptor de la app si ya aparece en la lista).
3. Selecciona el ejecutable de tu app (`.app` en Aplicaciones, o la ruta de build en desarrollo).
4. Activa el interruptor junto al nombre de la app.

**Nota de privacidad:** solo solicita este permiso si la funcionalidad está justificada en la política de privacidad del producto.

---

## 2. Automatización (Apple Events)

**Ruta:** *Privacidad y seguridad → Automatización*

**Qué permite:** que tu app envíe **AppleScript** u otros eventos de Apple a **otras aplicaciones** (Safari, Chrome, etc.) para leer estado (por ejemplo pestaña activa, URL).

**Cuándo lo necesita este proyecto:** estrategia recomendada para **Safari** y navegadores **Chromium** cuando existan diccionarios AppleScript estables.

**Cómo aparece:** la primera vez que el script controla otra app, macOS puede mostrar un diálogo pidiendo permiso. Luego la relación **tu app → app destino** puede listarse en *Automatización*.

**Si no aparece el diálogo:** ejecuta manualmente una acción que dispare el script; comprueba que la app destino esté instalada y permita automatización según su propia configuración.

---

## 3. Grabación de pantalla

**Ruta:** *Privacidad y seguridad → Grabación de pantalla*

**Qué permite:** capturar píxeles de la pantalla o ventanas.

**Cuándo considerarlo:** solo si el producto incluye una función explícita que lo requiera (p. ej. OCR de la barra de pestañas). **No** es necesario para el diseño base descrito en el brief (workspace + AppleScript + Accesibilidad).

---

## 4. Entrada completa (Input Monitoring)

**Ruta:** *Privacidad y seguridad → Entrada completa*

**Qué permite:** observar eventos de teclado a nivel del sistema.

**Cuándo lo necesita este proyecto:** en general **no**, salvo que implementes atajos globales que lo exijan.

---

## 5. Notificaciones

**Ruta:** *Ajustes del Sistema → Notificaciones → [tu app]*

**Uso:** avisos de fin de Pomodoro, recordatorios de foco, etc.

---

## 6. Elementos de inicio (Login Items)

**Ruta:** *Ajustes del Sistema → General → Elementos de inicio*

**Uso:** opcional, si quieres que un agente de tracking arranque al iniciar sesión.

---

## 7. Distribución y Gatekeeper

Para apps **firmadas con Developer ID** y **notarizadas**:

- Los usuarios suelen abrir la app sin pasos extra.
- Si macOS bloquea la primera ejecución: *Clic derecho en la app → Abrir* y confirmar.

Documenta en el README del producto cualquier paso específico de tu pipeline de firma.

---

## Resumen rápido (checklist)

| Permiso | Uso típico en este proyecto |
|---------|-----------------------------|
| Accesibilidad | Títulos de ventana / fallback de contexto |
| Automatización | Safari / Chromium vía AppleScript |
| Grabación de pantalla | Solo si implementas captura/OCR |
| Entrada completa | Solo atajos globales avanzados |
| Notificaciones | Pomodoro y avisos |

---

## Referencia

Brief del producto: [brief-app-tracking-pomodoro.md](brief-app-tracking-pomodoro.md)
