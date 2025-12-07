/* @Programa principal del proyecto SDGEStreaming - Programación orientada a objetos
   @Autores: Nelson Espinosa, Barbara Peñaherrera
   @Domingo 7 de diciembre de 2025. Quito - Ecuador
   @Punto de entrada del sistema. Contiene el menú interactivo y la lógica de control principal que orquesta las interacciones con los módulos y la base de datos.*/
// cmd/sdge/main.go
package main

import (
	"SDGEStreaming/internal/db"
	"SDGEStreaming/internal/models"
	"SDGEStreaming/internal/repositories"
	"SDGEStreaming/internal/security"
	"SDGEStreaming/internal/services"
	"SDGEStreaming/internal/utils"
	"fmt"
	"os"
	"time"
)

var (
	currentUser *CurrentUser
)

type CurrentUser struct {
	ID        int
	Name      string
	Email     string
	PlanID    int
	PlanName  string
	Age       int
	AgeRating string
	IsAdmin   bool
}

var (
	userService         *services.UserService
	contentService      *services.ContentService
	subscriptionService *services.SubscriptionService
	playbackService     *services.PlaybackService

	userRepo repositories.UserRepo
)

func main() {
	if err := db.InitDB("sdgestreaming.db"); err != nil {
		fmt.Printf("Error fatal al iniciar la base de datos: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	// INICIALIZAR REPOS
	userRepo = repositories.NewUserRepo()
	contentRepo := repositories.NewContentRepo()
	subscriptionRepo := repositories.NewSubscriptionRepo()
	playbackHistoryRepo := repositories.NewPlaybackHistoryRepo()
	favoriteRepo := repositories.NewFavoriteRepo()

	// Crear usuario admin si no existe
	adminUser, err := userRepo.FindByEmail("admin@sdge.com")
	if err != nil {
		fmt.Printf("Error buscando usuario admin: %v\n", err)
	}
	if adminUser == nil {
		hashedPass, err := security.HashPassword("admin123")
		if err != nil {
			fmt.Printf("Error generando contraseña del admin: %v\n", err)
		} else {
			now := time.Now()
			adminModel := &models.User{
				Name:         "Admin",
				Email:        "admin@sdge.com",
				Age:          30,
				PlanID:       3,
				AgeRating:    "Adulto",
				IsAdmin:      true,
				PasswordHash: hashedPass,
				CreatedAt:    now,
				LastLogin:    now,
			}
			if err := userRepo.Create(adminModel); err != nil {
				fmt.Printf("Error creando usuario admin: %v\n", err)
			}
		}
	}

	userService = services.NewUserService(userRepo, subscriptionRepo)
	contentService = services.NewContentService(contentRepo)
	subscriptionService = services.NewSubscriptionService(subscriptionRepo, userRepo)
	playbackService = services.NewPlaybackService(playbackHistoryRepo, favoriteRepo, contentRepo)

	utils.ClearScreen()
	runApplication()
}

func runApplication() {
	for {
		if currentUser == nil {
			showAuthMenu()
		} else {
			showMainMenu()
		}
	}
}

func showAuthMenu() {
	utils.ClearScreen()
	fmt.Println("╔════════════════════════════════╗")
	fmt.Println("║    SDGEStreaming - Inicio    ║")
	fmt.Println("╚════════════════════════════════╝")
	fmt.Println()
	fmt.Println("1. Iniciar Sesión")
	fmt.Println("2. Registrarse")
	fmt.Println("3. Salir")
	fmt.Print("\nSeleccione una opción: ")

	option := utils.ReadLine("")
	switch option {
	case "1":
		login()
	case "2":
		register()
	case "3":
		fmt.Println("¡Gracias por usar SDGEStreaming!")
		os.Exit(0)
	default:
		fmt.Println("Opción inválida.")
		utils.WaitForEnter()
	}
}

func login() {
	utils.ClearScreen()
	fmt.Println("Iniciar Sesión")
	fmt.Println("══════════════")
	email := utils.ReadLine("Email: ")
	password := utils.ReadLine("Contraseña: ")

	user, err := userService.Login(email, password)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		utils.WaitForEnter()
		return
	}

	planName := getPlanName(user.PlanID)

	currentUser = &CurrentUser{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		PlanID:    user.PlanID,
		PlanName:  planName,
		Age:       user.Age,
		AgeRating: user.AgeRating,
		IsAdmin:   user.IsAdmin,
	}
	fmt.Printf("¡Bienvenido, %s!\n", user.Name)
	utils.WaitForEnter()
}

func register() {
	utils.ClearScreen()
	fmt.Println("Registro de Nuevo Usuario")
	fmt.Println("══════════════════════════")
	name := utils.ReadLine("Nombre completo: ")
	ageStr := utils.ReadLine("Edad (mínimo 13): ")
	age, err := utils.ToInt(ageStr)
	if err != nil || age < 13 {
		fmt.Println("Edad inválida. Debe ser un número entero mayor o igual a 13.")
		utils.WaitForEnter()
		return
	}

	email := utils.ReadLine("Email: ")
	password := utils.ReadLine("Contraseña (mínimo 6 caracteres): ")

	if !utils.IsValidEmail(email) {
		fmt.Println("Formato de email inválido.")
		utils.WaitForEnter()
		return
	}
	if !utils.IsValidPassword(password) {
		fmt.Println("La contraseña debe tener al menos 6 caracteres.")
		utils.WaitForEnter()
		return
	}

	_, err = userService.Register(name, age, email, password, false)
	if err != nil {
		fmt.Printf("Error en el registro: %v\n", err)
	} else {
		fmt.Println("¡Registro exitoso! Ahora puede iniciar sesión.")
	}
	utils.WaitForEnter()
}

func showMainMenu() {
	utils.ClearScreen()
	fmt.Println("Menú Principal")
	fmt.Println("══════════════")
	fmt.Printf("Hola, %s (%s)\n", currentUser.Name, currentUser.PlanName)
	fmt.Println()
	fmt.Println("1. Inicio")
	fmt.Println("2. Tendencias")
	fmt.Println("3. Explorar Contenido")
	fmt.Println("4. Mi Lista")
	fmt.Println("5. Perfil y Cuenta")
	if currentUser.IsAdmin {
		fmt.Println("6. Panel de Administración")
		fmt.Println("7. Cerrar Sesión")
	} else {
		fmt.Println("6. Cerrar Sesión")
	}
	fmt.Print("\nSeleccione una opción: ")

	option := utils.ReadLine("")
	switch option {
	case "1":
		showHome()
	case "2":
		showTrending()
	case "3":
		browseContent(false)
	case "4":
		showMyList()
	case "5":
		showProfileMenu()
	case "6":
		if currentUser.IsAdmin {
			showAdminPanel()
		} else {
			logout()
		}
	case "7":
		if currentUser.IsAdmin {
			logout()
		}
	default:
		fmt.Println("Opción inválida.")
	}
	utils.WaitForEnter()
}

func showHome() {
	utils.ClearScreen()
	fmt.Println("Inicio")
	fmt.Println("══════")
	fmt.Println("¡Bienvenido a tu página de inicio!\n")

	fmt.Println("► Continuar viendo:")
	continueWatching, _ := playbackService.GetContinueWatching(currentUser.ID)
	if len(continueWatching) == 0 {
		fmt.Println("  No tienes nada en progreso.")
	} else {
		for _, entry := range continueWatching {
			var title string
			if entry.ContentType == "audiovisual" {
				content, _ := contentService.GetAudiovisualByID(entry.ContentID)
				if content != nil {
					title = content.Title
				}
			} else {
				content, _ := contentService.GetAudioByID(entry.ContentID)
				if content != nil {
					title = content.Title
				}
			}
			if title != "" {
				fmt.Printf("  * %s (ID: %d)\n", title, entry.ContentID)
			}
		}
	}

	fmt.Println()
	utils.WaitForEnter()
}

func showTrending() {
	utils.ClearScreen()
	fmt.Println("Tendencias")
	fmt.Println("══════════")

	fmt.Println("\n🎬 Contenido Audiovisual Popular:")
	audiovisuals, err := contentService.GetAllAudiovisual()
	if err != nil {
		fmt.Printf("Error al cargar contenido: %v\n", err)
	} else {
		for i, av := range audiovisuals {
			if i >= 3 {
				break
			}
			fmt.Printf("  %d. %s (%.1f⭐)\n", i+1, av.Title, av.AverageRating)
		}
	}

	fmt.Println("\n🎵 Contenido de Audio Popular:")
	audios, err := contentService.GetAllAudio()
	if err != nil {
		fmt.Printf("Error al cargar contenido: %v\n", err)
	} else {
		for i, a := range audios {
			if i >= 3 {
				break
			}
			fmt.Printf("  %d. %s - %s (%.1f⭐)\n", i+1, a.Artist, a.Title, a.AverageRating)
		}
	}

	utils.WaitForEnter()
}

func browseContent(isGuest bool) {
	for {
		utils.ClearScreen()
		fmt.Println("Explorar Contenido")
		fmt.Println("══════════════════")
		fmt.Println("1. Contenido Audiovisual")
		fmt.Println("2. Contenido de Audio")
		fmt.Println("3. Volver")
		fmt.Print("\nSeleccione una opción: ")

		option := utils.ReadLine("")
		switch option {
		case "1":
			browseAudiovisual(isGuest)
		case "2":
			browseAudio(isGuest)
		case "3":
			return
		default:
			fmt.Println("Opción inválida.")
			utils.WaitForEnter()
		}
	}
}

func browseAudiovisual(isGuest bool) {
	contents, err := contentService.GetAllAudiovisualForUser(currentUser.AgeRating)
	if err != nil {
		fmt.Printf("Error al cargar contenido: %v\n", err)
		utils.WaitForEnter()
		return
	}

	if len(contents) == 0 {
		fmt.Println("No hay contenido audiovisual disponible para tu clasificación de edad.")
		utils.WaitForEnter()
		return
	}

	utils.ClearScreen()
	fmt.Println("\n🎬 Contenido Audiovisual Disponible:")
	for _, c := range contents {
		fmt.Printf("ID: %d | %s (%s)\n", c.ID, c.Title, c.Type)
		fmt.Printf("   Género: %s | Duración: %d min | Clasificación: %s\n", c.Genre, c.Duration, c.AgeRating)
		fmt.Printf("   Promedio: %.1f⭐\n", c.AverageRating)
		fmt.Println("────────────────────────────────────────")
	}

	if !isGuest {
		contentIDStr := utils.ReadLine("\nIngrese el ID del contenido para ver detalles (0 para volver): ")
		if contentIDStr == "0" {
			return
		}
		contentID, err := utils.ToInt(contentIDStr)
		if err != nil {
			fmt.Println("ID inválido.")
			utils.WaitForEnter()
			return
		}

		content, err := contentService.GetAudiovisualByID(contentID)
		if err != nil {
			fmt.Println("Contenido no encontrado.")
			utils.WaitForEnter()
			return
		}

		utils.ClearScreen()
		fmt.Printf("═══ %s ═══\n", content.Title)
		fmt.Printf("Tipo: %s\n", content.Type)
		fmt.Printf("Género: %s\n", content.Genre)
		fmt.Printf("Sinopsis: %s\n", content.Synopsis)
		fmt.Printf("Director: %s\n", content.Director)
		fmt.Printf("Año: %d\n", content.ReleaseYear)
		fmt.Printf("Duración: %d minutos\n", content.Duration)
		fmt.Printf("Clasificación: %s\n", content.AgeRating)
		fmt.Printf("Promedio de calificación: %.1f⭐\n", content.AverageRating)
		fmt.Println("\n1. Reproducir")
		fmt.Println("2. Marcar como favorito")
		fmt.Println("3. Calificar")
		fmt.Println("4. Volver")
		action := utils.ReadLine("Seleccione una acción: ")

		switch action {
		case "1":
			playAudiovisual(contentID)
		case "2":
			err = playbackService.AddFavorite(currentUser.ID, contentID, "audiovisual")
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			} else {
				fmt.Println("¡Agregado a Mi Lista!")
			}
			utils.WaitForEnter()
		case "3":
			rateContent(contentID, "audiovisual")
		case "4":
			return
		}
	}
}

func browseAudio(isGuest bool) {
	contents, err := contentService.GetAllAudioForUser(currentUser.AgeRating)
	if err != nil {
		fmt.Printf("Error al cargar contenido: %v\n", err)
		utils.WaitForEnter()
		return
	}

	if len(contents) == 0 {
		fmt.Println("No hay contenido de audio disponible para tu clasificación de edad.")
		utils.WaitForEnter()
		return
	}

	utils.ClearScreen()
	fmt.Println("\n🎵 Contenido de Audio Disponible:")
	for _, c := range contents {
		fmt.Printf("ID: %d | %s - %s\n", c.ID, c.Artist, c.Title)
		fmt.Printf("   Tipo: %s | Género: %s | Álbum: %s\n", c.Type, c.Genre, c.Album)
		fmt.Printf("   Duración: %d min | Clasificación: %s\n", c.Duration, c.AgeRating)
		fmt.Printf("   Promedio: %.1f⭐\n", c.AverageRating)
		fmt.Println("────────────────────────────────────────")
	}

	if !isGuest {
		contentIDStr := utils.ReadLine("\nIngrese el ID del contenido para ver detalles (0 para volver): ")
		if contentIDStr == "0" {
			return
		}
		contentID, err := utils.ToInt(contentIDStr)
		if err != nil {
			fmt.Println("ID inválido.")
			utils.WaitForEnter()
			return
		}

		content, err := contentService.GetAudioByID(contentID)
		if err != nil {
			fmt.Println("Contenido no encontrado.")
			utils.WaitForEnter()
			return
		}

		utils.ClearScreen()
		fmt.Printf("═══ %s ═══\n", content.Title)
		fmt.Printf("Artista: %s\n", content.Artist)
		fmt.Printf("Álbum: %s\n", content.Album)
		fmt.Printf("Género: %s\n", content.Genre)
		fmt.Printf("Duración: %d minutos\n", content.Duration)
		fmt.Printf("Clasificación: %s\n", content.AgeRating)
		fmt.Printf("Promedio de calificación: %.1f⭐\n", content.AverageRating)
		fmt.Println("\n1. Reproducir")
		fmt.Println("2. Marcar como favorito")
		fmt.Println("3. Calificar")
		fmt.Println("4. Volver")
		action := utils.ReadLine("Seleccione una acción: ")

		switch action {
		case "1":
			playAudio(contentID)
		case "2":
			err = playbackService.AddFavorite(currentUser.ID, contentID, "audio")
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			} else {
				fmt.Println("¡Agregado a Mi Lista!")
			}
			utils.WaitForEnter()
		case "3":
			rateContent(contentID, "audio")
		case "4":
			return
		}
	}
}

func showMyList() {
	utils.ClearScreen()
	fmt.Println("Mi Lista")
	fmt.Println("════════")

	favorites, err := playbackService.GetFavorites(currentUser.ID)
	if err != nil {
		fmt.Printf("Error al cargar favoritos: %v\n", err)
		utils.WaitForEnter()
		return
	}

	if len(favorites) == 0 {
		fmt.Println("No tienes ningún contenido en tu lista.")
		utils.WaitForEnter()
		return
	}

	fmt.Println("Contenido en tu lista:")
	for _, fav := range favorites {
		var title, details string
		if fav.ContentType == "audiovisual" {
			content, _ := contentService.GetAudiovisualByID(fav.ContentID)
			if content != nil {
				title = content.Title
				details = fmt.Sprintf("[%s] %s", content.Type, content.Genre)
			}
		} else {
			content, _ := contentService.GetAudioByID(fav.ContentID)
			if content != nil {
				title = fmt.Sprintf("%s - %s", content.Artist, content.Title)
				details = fmt.Sprintf("[%s] %s", content.Type, content.Genre)
			}
		}
		if title != "" {
			fmt.Printf("* %s\n", title)
			fmt.Printf("  %s\n", details)
			fmt.Println("────────────────────────────────────────")
		}
	}
	utils.WaitForEnter()
}

func showProfileMenu() {
	for {
		utils.ClearScreen()
		fmt.Println("Mi Perfil")
		fmt.Println("═════════")
		fmt.Printf("Nombre: %s\n", currentUser.Name)
		fmt.Printf("Email: %s\n", currentUser.Email)
		fmt.Printf("Plan actual: %s\n", currentUser.PlanName)
		fmt.Printf("Edad: %d\n", currentUser.Age)
		fmt.Printf("Clasificación: %s\n", currentUser.AgeRating)
		fmt.Println()
		fmt.Println("1. Cambiar Plan de Suscripción")
		fmt.Println("2. Ver Métodos de Pago")
		fmt.Println("3. Ver Historial de Reproducción")
		fmt.Println("4. Volver al Menú Principal")
		fmt.Print("\nSeleccione una opción: ")

		option := utils.ReadLine("")
		switch option {
		case "1":
			upgradePlan()
		case "2":
			viewPaymentMethods()
		case "3":
			viewPlaybackHistory()
		case "4":
			return
		default:
			fmt.Println("Opción inválida.")
			utils.WaitForEnter()
		}
	}
}

func viewPlaybackHistory() {
	utils.ClearScreen()
	fmt.Println("Historial de Reproducción")
	fmt.Println("══════════════════════════")

	history, err := playbackService.GetHistory(currentUser.ID)
	if err != nil {
		fmt.Printf("Error al cargar el historial: %v\n", err)
		utils.WaitForEnter()
		return
	}

	if len(history) == 0 {
		fmt.Println("No tienes historial de reproducción.")
		utils.WaitForEnter()
		return
	}

	fmt.Println("Tus últimas reproducciones:")
	for _, entry := range history {
		var title string
		if entry.ContentType == "audiovisual" {
			content, _ := contentService.GetAudiovisualByID(entry.ContentID)
			if content != nil {
				title = content.Title
			}
		} else {
			content, _ := contentService.GetAudioByID(entry.ContentID)
			if content != nil {
				title = content.Title
			}
		}
		if title != "" {
			fmt.Printf("* %s (%s)\n", title, entry.ContentType)
		}
	}
	utils.WaitForEnter()
}

func upgradePlan() {
	utils.ClearScreen()
	fmt.Println("Cambiar Plan de Suscripción")
	fmt.Println("════════════════════════════")

	plans, err := subscriptionService.GetAvailablePlans()
	if err != nil {
		fmt.Printf("Error al cargar planes: %v\n", err)
		utils.WaitForEnter()
		return
	}

	fmt.Println("Planes disponibles:")
	for _, p := range plans {
		fmt.Printf("%d. %s - $%.2f/mes | Calidad: %s | Dispositivos: %d\n",
			p.ID, p.Name, p.Price, p.MaxQuality, p.MaxDevices)
	}

	planIDStr := utils.ReadLine("\nSeleccione el número del plan deseado (0 para cancelar): ")
	if planIDStr == "0" {
		return
	}
	planID, err := utils.ToInt(planIDStr)
	if err != nil {
		fmt.Println("Selección inválida.")
		utils.WaitForEnter()
		return
	}

	if planID == currentUser.PlanID {
		fmt.Println("Ya está suscrito a este plan.")
		utils.WaitForEnter()
		return
	}

	if planID == 1 {
		err = userService.UpdateUserPlan(currentUser.ID, 1)
		if err != nil {
			fmt.Printf("Error al actualizar el plan: %v\n", err)
			utils.WaitForEnter()
			return
		}
		currentUser.PlanID = 1
		currentUser.PlanName = "Free"
		fmt.Println("Su plan ha sido cambiado a Free.")
		utils.WaitForEnter()
		return
	}

	fmt.Println("\n--- Información de Pago ---")
	cardHolder := utils.ReadLine("Nombre del titular de la tarjeta: ")
	cardNumber := utils.ReadLine("Número de tarjeta (16 dígitos): ")
	expiry := utils.ReadLine("Fecha de vencimiento (MM/AAAA): ")
	cvvStr := utils.ReadLine("CVV (3 dígitos): ")

	var expiryMonth, expiryYear int
	if len(expiry) == 7 && expiry[2] == '/' {
		expiryMonth, _ = utils.ToInt(expiry[0:2])
		expiryYear, _ = utils.ToInt(expiry[3:7])
	} else {
		fmt.Println("Formato de fecha de vencimiento inválido (MM/AAAA).")
		utils.WaitForEnter()
		return
	}

	cvv, err := utils.ToInt(cvvStr)
	if err != nil {
		fmt.Println("CVV inválido.")
		utils.WaitForEnter()
		return
	}

	err = subscriptionService.ProcessPayment(currentUser.ID, planID, cardHolder, cardNumber, expiryMonth, expiryYear, cvv)
	if err != nil {
		fmt.Printf("Error en el pago: %v\n", err)
		utils.WaitForEnter()
		return
	}

	currentUser.PlanID = planID
	currentUser.PlanName = getPlanName(planID)
	fmt.Println("¡Su plan ha sido actualizado exitosamente!")
	utils.WaitForEnter()
}

func viewPaymentMethods() {
	utils.ClearScreen()
	fmt.Println("Métodos de Pago")
	fmt.Println("════════════════")
	method, err := userService.GetDefaultPaymentMethod(currentUser.ID)
	if err != nil {
		fmt.Println("No tiene métodos de pago guardados.")
	} else {
		fmt.Printf("Tarjeta predeterminada: **** **** **** %s\n", method.Last4)
		fmt.Printf("Titular: %s\n", method.CardHolder)
		fmt.Printf("Vence: %02d/%d\n", method.ExpiryMonth, method.ExpiryYear)
	}
	utils.WaitForEnter()
}

func showAdminPanel() {
	utils.ClearScreen()
	fmt.Println("Panel de Administración")
	fmt.Println("════════════════════════")
	fmt.Println("1. Gestionar Usuarios")
	fmt.Println("2. Gestionar Contenido")
	fmt.Println("3. Generar Reportes")
	fmt.Println("4. Volver")
	fmt.Print("\nSeleccione una opción: ")

	option := utils.ReadLine("")
	switch option {
	case "1":
		manageUsers()
	case "2":
		manageContent()
	case "3":
		generateReports()
	case "4":
		return
	default:
		fmt.Println("Opción inválida.")
	}
	utils.WaitForEnter()
}

func manageUsers() {
	utils.ClearScreen()
	fmt.Println("Gestión de Usuarios")
	fmt.Println("════════════════════")
	users, err := userService.GetAllUsers()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		utils.WaitForEnter()
		return
	}
	for _, u := range users {
		adminTag := ""
		if u.IsAdmin {
			adminTag = " [ADMIN]"
		}
		fmt.Printf("- %s%s (%s) | Edad: %d | Clasificación: %s\n", u.Name, adminTag, u.Email, u.Age, u.AgeRating)
	}
	utils.WaitForEnter()
}

func manageContent() {
	for {
		utils.ClearScreen()
		fmt.Println("Gestión de Contenido")
		fmt.Println("════════════════════")
		fmt.Println("1. Agregar Contenido Audiovisual")
		fmt.Println("2. Agregar Contenido de Audio")
		fmt.Println("3. Listar Contenido Audiovisual")
		fmt.Println("4. Listar Contenido de Audio")
		fmt.Println("5. Volver")
		fmt.Print("\nSeleccione una opción: ")

		option := utils.ReadLine("")
		switch option {
		case "1":
			addAudiovisualContent()
		case "2":
			addAudioContent()
		case "3":
			listAudiovisualAdmin()
		case "4":
			listAudioAdmin()
		case "5":
			return
		default:
			fmt.Println("Opción inválida.")
			utils.WaitForEnter()
		}
	}
}

func addAudiovisualContent() {
	utils.ClearScreen()
	fmt.Println("Agregar Contenido Audiovisual")
	fmt.Println("══════════════════════════════")

	title := utils.ReadLine("Título: ")
	contentType := utils.ReadLine("Tipo (movie/series/documentary): ")
	genre := utils.ReadLine("Género: ")
	durationStr := utils.ReadLine("Duración (minutos): ")
	duration, err := utils.ToInt(durationStr)
	if err != nil {
		fmt.Println("Duración inválida.")
		utils.WaitForEnter()
		return
	}

	ageRating := utils.ReadLine("Clasificación (G/PG/PG-13/R): ")
	synopsis := utils.ReadLine("Sinopsis: ")
	yearStr := utils.ReadLine("Año de lanzamiento: ")
	year, err := utils.ToInt(yearStr)
	if err != nil {
		fmt.Println("Año inválido.")
		utils.WaitForEnter()
		return
	}
	director := utils.ReadLine("Director: ")

	err = contentService.CreateAudiovisual(title, contentType, genre, duration, ageRating, synopsis, year, director)
	if err != nil {
		fmt.Printf("Error al agregar contenido: %v\n", err)
	} else {
		fmt.Println("¡Contenido agregado exitosamente!")
	}
	utils.WaitForEnter()
}

func addAudioContent() {
	utils.ClearScreen()
	fmt.Println("Agregar Contenido de Audio")
	fmt.Println("═══════════════════════════")

	title := utils.ReadLine("Título: ")
	contentType := utils.ReadLine("Tipo (song/podcast/audiobook): ")
	genre := utils.ReadLine("Género: ")
	durationStr := utils.ReadLine("Duración (minutos): ")
	duration, err := utils.ToInt(durationStr)
	if err != nil {
		fmt.Println("Duración inválida.")
		utils.WaitForEnter()
		return
	}

	ageRating := utils.ReadLine("Clasificación (General/Explicit): ")
	artist := utils.ReadLine("Artista: ")
	album := utils.ReadLine("Álbum: ")
	trackStr := utils.ReadLine("Número de pista: ")
	trackNumber, err := utils.ToInt(trackStr)
	if err != nil {
		trackNumber = 1
	}

	err = contentService.CreateAudio(title, contentType, genre, duration, ageRating, artist, album, trackNumber)
	if err != nil {
		fmt.Printf("Error al agregar contenido: %v\n", err)
	} else {
		fmt.Println("¡Contenido agregado exitosamente!")
	}
	utils.WaitForEnter()
}

func listAudiovisualAdmin() {
	utils.ClearScreen()
	fmt.Println("Lista de Contenido Audiovisual")
	fmt.Println("═══════════════════════════════")

	contents, err := contentService.GetAllAudiovisual()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		utils.WaitForEnter()
		return
	}

	for _, c := range contents {
		fmt.Printf("ID: %d | %s (%s) - %d min | Clasificación: %s\n", c.ID, c.Title, c.Type, c.Duration, c.AgeRating)
	}
	utils.WaitForEnter()
}

func listAudioAdmin() {
	utils.ClearScreen()
	fmt.Println("Lista de Contenido de Audio")
	fmt.Println("════════════════════════════")

	contents, err := contentService.GetAllAudio()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		utils.WaitForEnter()
		return
	}

	for _, c := range contents {
		fmt.Printf("ID: %d | %s - %s (%s) - %d min | Clasificación: %s\n", c.ID, c.Artist, c.Title, c.Type, c.Duration, c.AgeRating)
	}
	utils.WaitForEnter()
}

func generateReports() {
	utils.ClearScreen()
	fmt.Println("Generación de Reportes")
	fmt.Println("═══════════════════════")
	users, err := userService.GetAllUsers()
	if err != nil {
		fmt.Printf("Error al cargar usuarios: %v\n", err)
	} else {
		fmt.Printf("* Total de Usuarios: %d\n", len(users))
	}

	audiovisuals, err := contentService.GetAllAudiovisual()
	if err != nil {
		fmt.Printf("Error al cargar contenido audiovisual: %v\n", err)
	} else {
		fmt.Printf("* Total de Contenido Audiovisual: %d\n", len(audiovisuals))
		if len(audiovisuals) > 0 {
			fmt.Printf("* Contenido más popular: '%s' (%.1f⭐)\n", audiovisuals[0].Title, audiovisuals[0].AverageRating)
		}
	}

	audios, err := contentService.GetAllAudio()
	if err != nil {
		fmt.Printf("Error al cargar contenido de audio: %v\n", err)
	} else {
		fmt.Printf("* Total de Contenido de Audio: %d\n", len(audios))
		if len(audios) > 0 {
			fmt.Printf("* Audio más popular: '%s - %s' (%.1f⭐)\n", audios[0].Artist, audios[0].Title, audios[0].AverageRating)
		}
	}

	utils.WaitForEnter()
}

func logout() {
	currentUser = nil
	fmt.Println("Sesión cerrada correctamente.")
	utils.WaitForEnter()
}

func getPlanName(planID int) string {
	switch planID {
	case 1:
		return "Free"
	case 2:
		return "Estándar"
	case 3:
		return "Premium 4K"
	default:
		return "Desconocido"
	}
}

func playAudiovisual(contentID int) {
	content, err := contentService.GetAudiovisualByID(contentID)
	if err != nil {
		fmt.Println("Error al cargar contenido.")
		utils.WaitForEnter()
		return
	}

	utils.ClearScreen()
	fmt.Printf("▶ Reproduciendo: %s\n", content.Title)
	fmt.Println("══════════════════════════════════════")
	fmt.Println("Simulando reproducción...")
	fmt.Println("[████████████████████] 100%")
	fmt.Printf("Duración total: %d minutos\n", content.Duration)
	fmt.Println("══════════════════════════════════════")

	// Registrar en historial
	if err := playbackService.AddToHistory(currentUser.ID, contentID, "audiovisual"); err != nil {
		fmt.Printf("No se pudo registrar en historial: %v\n", err)
	}

	// Simular progreso (50% visto)
	progressSeconds := (content.Duration * 60) / 2
	if err := playbackService.UpdateProgress(currentUser.ID, contentID, "audiovisual", progressSeconds); err != nil {
		fmt.Printf("No se pudo actualizar progreso: %v\n", err)
	}

	fmt.Println("\n✓ Reproducción finalizada")
	fmt.Println("Se ha guardado tu progreso.")
	utils.WaitForEnter()
}

func playAudio(contentID int) {
	content, err := contentService.GetAudioByID(contentID)
	if err != nil {
		fmt.Println("Error al cargar contenido.")
		utils.WaitForEnter()
		return
	}

	utils.ClearScreen()
	fmt.Printf("♪ Reproduciendo: %s - %s\n", content.Artist, content.Title)
	fmt.Println("══════════════════════════════════════")
	fmt.Println("Simulando reproducción...")
	fmt.Println("[████████████████████] 100%")
	fmt.Printf("Duración total: %d minutos\n", content.Duration)
	fmt.Println("══════════════════════════════════════")

	// Registrar en historial
	if err := playbackService.AddToHistory(currentUser.ID, contentID, "audio"); err != nil {
		fmt.Printf("No se pudo registrar en historial: %v\n", err)
	}

	// Simular progreso (70% escuchado)
	progressSeconds := (content.Duration * 60) * 7 / 10
	if err := playbackService.UpdateProgress(currentUser.ID, contentID, "audio", progressSeconds); err != nil {
		fmt.Printf("No se pudo actualizar progreso: %v\n", err)
	}

	fmt.Println("\n✓ Reproducción finalizada")
	fmt.Println("Se ha guardado tu progreso.")
	utils.WaitForEnter()
}

func rateContent(contentID int, contentType string) {
	utils.ClearScreen()
	fmt.Println("Calificar Contenido")
	fmt.Println("════════════════════")
	fmt.Println("Ingrese su calificación (1.0 - 10.0)")
	ratingStr := utils.ReadLine("Calificación: ")

	rating, err := utils.ToFloat(ratingStr)
	if err != nil || rating < 1.0 || rating > 10.0 {
		fmt.Println("Calificación inválida. Debe ser entre 1.0 y 10.0")
		utils.WaitForEnter()
		return
	}

	// Guardar calificación
	err = contentService.RateContent(currentUser.ID, contentID, contentType, rating)
	if err != nil {
		fmt.Printf("Error al calificar: %v\n", err)
	} else {
		fmt.Printf("¡Gracias! Has calificado este contenido con %.1f⭐\n", rating)
	}
	utils.WaitForEnter()
}
