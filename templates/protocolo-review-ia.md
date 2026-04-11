# Checklist de Revisión

## 1. Alucinaciones de librerías

### Pregunta clave

- [ ] Los imports existen
- [ ] Las funciones usadas son reales
- [ ] La API corresponde a la versión actual
- [ ] Existen vulnerabilidades en estas versiones

## 2. Lógica de negocio sutil

- [ ] Los cálculos son correctos
- [ ] El manejo de fechas es consistente y en UTF-8
- [ ] No se usa float para dinero
- [ ] Edge cases están contemplados

## 3. Seguridad

### Pregunta clave

- [ ] Inputs están validados
- [ ] No hay riesgo de inyección
- [ ] No se exponen credenciales
- [ ] No se filtran datos sensibles

## 4. Context window

- [ ] El código respeta el brief original
- [ ] Se cumplen los constraints definidos
- [ ] No se agregaron dependencias innecesarias
- [ ] Se cumple la Definition of Done

## 5. Punto personalizado del proyecto

- [ ] Tests cubren el código nuevo
- [ ] Logs no exponen datos sensibles
- [ ] Performance es aceptable
- [ ] Métricas funcionan correctamente
- [ ] No acoplar elementos de la arquitectura en el momento que se deba usar un proceso invocar el servicio no el code.

## Resultado del Review

- [ ] El código pasa los 5 puntos del protocolo
- [ ] Se hicieron correcciones necesarias
- [ ] El código está listo para commit

**Reviewer:** ____________________
**Fecha:** ____________________

---