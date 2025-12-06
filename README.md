
> **No se usan librerías externas.** Todo se implementó con la biblioteca estándar de Go.

---

## Funcionalidades Implementadas (AA1)

| Funcionalidad | Descripción |
|---------------|-------------|
| **Registro de usuarios** | Validación de email y contraseña (mínimo 6 caracteres). |
| **Inicio de sesión** | Autenticación por email y contraseña. Usuario administrador predeterminado: `admin@sdge.com / admin123`. |
| **Explorar contenido** | Catálogo de películas, series, música y podcasts con duración, género y clasificación por edad. |
| **Clasificación por edad** | Bloqueo automático de contenido no adecuado para la edad del usuario. |
| **Calificar contenido** | Dar calificación de 1.0 a 10.0. Se permite sobrescribir calificaciones anteriores con mensaje de confirmación. |
| **Promedios automáticos** | El sistema recalcula el rating promedio cada vez que se califica. |
| **Menús jerárquicos** | Navegación intuitiva con opción “0” para volver atrás en cualquier menú. |
| **Gestión de administrador** | Listar usuarios, agregar contenido audiovisual o de audio. |
| **Manejo de errores** | Mensajes claros y útiles. El programa no se cierra por entradas inválidas. |
| **Interfaz limpia** | Salida en consola con formato ordenado, sin colores ni dependencias externas. |

---

## Cómo Ejecutar el Proyecto

1. Asegúrate de tener **Go instalado** (versión 1.20 o superior).
2. Clona o descarga el proyecto:
   ```bash
   git clone https://github.com/tuusuario/SDGEStreaming.git
   cd SDGEStreaming
   go mod init SDGEStreaming
   go run cmd/sdge/main.go