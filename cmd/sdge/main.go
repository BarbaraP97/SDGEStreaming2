/* @Programa principal del proyecto SDGEStreaming - Programación orientada a objetos
   @Autores: Nelson Espinosa, Barbara Peñaherrera
   @Domingo 7 de diciembre de 2025. Quito - Ecuador
   @Punto de entrada del sistema. Contiene el menú interactivo y la lógica de control principal que orquesta las interacciones con los módulos y la base de datos.*/
// cmd/sdge/main.go
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
	"strings"
	"time"
)

var (
	currentUser    *models.User
	currentProfile *models.Profile
)

var (
	userService         *services.UserService
	profileService      *services.ProfileService
	contentService      *services.ContentService
	subscriptionService *services.SubscriptionService
	playbackService     *services.PlaybackService
	playlistService     *services.PlaylistService
	reportService       *services.ReportService
)

func main() {
	if err := db.InitDB("sdgestreaming.db"); err != nil {
		fmt.Printf("[ERROR] Error fatal al iniciar la base de datos: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	// Inicializar repositorios
	userRepo := repositories.NewUserRepo()
	profileRepo := repositories.NewProfileRepo()
	contentRepo := repositories.NewContentRepo()
	subscriptionRepo := repositories.NewSubscriptionRepo()
	playbackHistoryRepo := repositories.NewPlaybackHistoryRepo()
	favoriteRepo := repositories.NewFavoriteRepo()
	playlistRepo := repositories.NewPlaylistRepo()

	// Inicializar servicios
	userService = services.NewUserService(userRepo, subscriptionRepo)
	profileService = services.NewProfileService(profileRepo, userRepo)
	contentService = services.NewContentService(contentRepo)
	subscriptionService = services.NewSubscriptionService(subscriptionRepo, userRepo)
	playbackService = services.NewPlaybackService(playbackHistoryRepo, favoriteRepo, contentRepo)
	playlistService = services.NewPlaylistService(playlistRepo, contentRepo)
	reportService = services.NewReportService(userRepo, contentRepo, playbackHistoryRepo, subscriptionRepo)

	// Crear usuario admin si no existe
	adminUser, _ := userRepo.FindByEmail("admin@sdge.com")
	if adminUser == nil {
		hashedPass, _ := security.HashPassword("admin123")
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
		userRepo.Create(adminModel)
		
		// Crear perfil principal para admin
		profileService.CreateProfile(adminModel.ID, "Admin", "adult", true)
	}

	// Contenido de ejemplo
	initSampleContent()

	utils.ClearScreen()
	runApplication()
}

func runApplication() {
	for {
		if currentUser == nil {
			showAuthMenu()
		} else if currentProfile == nil {
			showProfileSelection()
		} else {
			showMainMenu()
		}
	}
}

func showAuthMenu() {
	utils.ClearScreen()
	fmt.Println("=========================================")
	fmt.Println("       SDGEStreaming - Bienvenido       ")
	fmt.Println("=========================================")
	fmt.Println()
	fmt.Println("  [1] Iniciar Sesion")
	fmt.Println("  [2] Registrarse")
	fmt.Println("  [0] Salir")
	fmt.Println()
	fmt.Print("Seleccione una opcion: ")

	option := utils.ReadLine("")
	switch option {
	case "1":
		login()
	case "2":
		register()
	case "0":
		fmt.Println("\nGracias por usar SDGEStreaming!")
		os.Exit(0)
	default:
		fmt.Println("\n[ERROR] Opcion invalida")
		time.Sleep(1 * time.Second)
	}
}

func login() {
	utils.ClearScreen()
	fmt.Println("=========================================")
	fmt.Println("            Iniciar Sesion              ")
	fmt.Println("=========================================")
	fmt.Println()
	email := utils.ReadLine("Email: ")
	password := utils.ReadLine("Contrasena: ")

	user, err := userService.Login(email, password)
	if err != nil {
		fmt.Printf("\n[ERROR] %v\n", err)
		time.Sleep(2 * time.Second)
		return
	}

	currentUser = user
	fmt.Printf("\n[OK] Bienvenido, %s!\n", user.Name)
	time.Sleep(1 * time.Second)
}

func register() {
	utils.ClearScreen()
	fmt.Println("=========================================")
	fmt.Println("        Registro de Nuevo Usuario       ")
	fmt.Println("=========================================")
	fmt.Println()
	
	name := utils.ReadLine("Nombre completo: ")
	ageStr := utils.ReadLine("Edad (minimo 13): ")
	age, err := utils.ToInt(ageStr)
	if err != nil || age < 13 {
		fmt.Println("\n[ERROR] Edad invalida. Debe ser mayor o igual a 13.")
		time.Sleep(2 * time.Second)
		return
	}

	email := utils.ReadLine("Email: ")
	password := utils.ReadLine("Contrasena (minimo 6 caracteres): ")

	if !utils.IsValidEmail(email) {
		fmt.Println("\n[ERROR] Formato de email invalido.")
		time.Sleep(2 * time.Second)
		return
	}
	if !utils.IsValidPassword(password) {
		fmt.Println("\n[ERROR] La contrasena debe tener al menos 6 caracteres.")
		time.Sleep(2 * time.Second)
		return
	}

	user, err := userService.Register(name, age, email, password, "Adulto", false)
	if err != nil {
		fmt.Printf("\n[ERROR] Error en el registro: %v\n", err)
		time.Sleep(2 * time.Second)
	} else {
		// Crear perfil principal
		profileService.CreateProfile(user.ID, name, "adult", true)
		fmt.Println("\n[OK] Registro exitoso! Ahora puede iniciar sesion.")
		time.Sleep(2 * time.Second)
	}
}

func showProfileSelection() {
	utils.ClearScreen()
	profiles, err := profileService.GetProfilesByUserID(currentUser.ID)
	if err != nil || len(profiles) == 0 {
		fmt.Println("[ERROR] No se encontraron perfiles")
		currentUser = nil
		return
	}

	fmt.Println("=========================================")
	fmt.Println("         Seleccionar Perfil             ")
	fmt.Println("=========================================")
	fmt.Println()
	for i, p := range profiles {
		profileType := map[string]string{"kids": "Ninos", "teen": "Joven", "adult": "Adulto"}[p.Type]
		fmt.Printf("  [%d] %s (%s)\n", i+1, p.Name, profileType)
	}
	fmt.Println()
	fmt.Println("  [A] Administrar Perfiles")
	fmt.Println("  [0] Cerrar Sesion")
	fmt.Println()
	fmt.Print("Seleccione: ")

	option := utils.ReadLine("")
	if option == "0" {
		currentUser = nil
		return
	}
	if strings.ToUpper(option) == "A" {
		manageProfiles()
		return
	}

	idx, err := utils.ToInt(option)
	if err != nil || idx < 1 || idx > len(profiles) {
		fmt.Println("\n[ERROR] Seleccion invalida")
		time.Sleep(1 * time.Second)
		return
	}

	currentProfile = &profiles[idx-1]
	fmt.Printf("\n[OK] Perfil seleccionado: %s\n", currentProfile.Name)
	time.Sleep(1 * time.Second)
}

func manageProfiles() {
	for {
		utils.ClearScreen()
		profiles, _ := profileService.GetProfilesByUserID(currentUser.ID)
		
		fmt.Println("=========================================")
		fmt.Println("         Administrar Perfiles           ")
		fmt.Println("=========================================")
		fmt.Println()
		
		for i, p := range profiles {
			profileType := map[string]string{"kids": "Ninos", "teen": "Joven", "adult": "Adulto"}[p.Type]
			mainTag := ""
			if p.IsMain {
				mainTag = " [PRINCIPAL]"
			}
			fmt.Printf("  [%d] %s (%s)%s\n", i+1, p.Name, profileType, mainTag)
		}
		
		fmt.Println()
		fmt.Println("  [A] Agregar Perfil")
		fmt.Println("  [E] Editar Perfil")
		fmt.Println("  [D] Eliminar Perfil")
		fmt.Println("  [0] Volver")
		fmt.Println()
		fmt.Print("Seleccione: ")

		option := strings.ToUpper(utils.ReadLine(""))
		switch option {
		case "A":
			addProfile()
		case "E":
			editProfile(profiles)
		case "D":
			deleteProfile(profiles)
		case "0":
			return
		default:
			fmt.Println("\n[ERROR] Opcion invalida")
			time.Sleep(1 * time.Second)
		}
	}
}

func addProfile() {
	utils.ClearScreen()
	fmt.Println("=========================================")
	fmt.Println("           Agregar Perfil               ")
	fmt.Println("=========================================")
	fmt.Println()
	
	count, _ := profileService.CountProfilesByUserID(currentUser.ID)
	maxProfiles := getMaxProfilesForPlan(currentUser.PlanID)
	
	if count >= maxProfiles {
		fmt.Printf("\n[ERROR] Ha alcanzado el limite de perfiles para su plan (%d)\n", maxProfiles)
		fmt.Println("Presione Enter para continuar...")
		utils.ReadLine("")
		return
	}

	name := utils.ReadLine("Nombre del perfil: ")
	fmt.Println("\nTipo de perfil:")
	fmt.Println("  [1] Ninos (contenido G)")
	fmt.Println("  [2] Joven (contenido G, PG, PG-13)")
	fmt.Println("  [3] Adulto (todo el contenido)")
	typeOption := utils.ReadLine("Seleccione: ")

	profileType := "adult"
	switch typeOption {
	case "1":
		profileType = "kids"
	case "2":
		profileType = "teen"
	case "3":
		profileType = "adult"
	}

	_, err := profileService.CreateProfile(currentUser.ID, name, profileType, false)
	if err != nil {
		fmt.Printf("\n[ERROR] %v\n", err)
	} else {
		fmt.Println("\n[OK] Perfil creado exitosamente")
	}
	time.Sleep(2 * time.Second)
}

func editProfile(profiles []models.Profile) {
	fmt.Print("\nIngrese el numero del perfil a editar (0 para cancelar): ")
	idx, err := utils.ToInt(utils.ReadLine(""))
	if err != nil || idx < 1 || idx > len(profiles) {
		return
	}

	profile := &profiles[idx-1]
	utils.ClearScreen()
	fmt.Println("=========================================")
	fmt.Printf("         Editar Perfil: %s\n", profile.Name)
	fmt.Println("=========================================")
	fmt.Println()

	newName := utils.ReadLine(fmt.Sprintf("Nuevo nombre [%s]: ", profile.Name))
	if newName != "" {
		profile.Name = newName
	}

	fmt.Println("\nTipo de perfil:")
	fmt.Println("  [1] Ninos")
	fmt.Println("  [2] Joven")
	fmt.Println("  [3] Adulto")
	typeOption := utils.ReadLine("Seleccione (Enter para mantener): ")

	switch typeOption {
	case "1":
		profile.Type = "kids"
		profile.AgeRating = "G"
	case "2":
		profile.Type = "teen"
		profile.AgeRating = "PG-13"
	case "3":
		profile.Type = "adult"
		profile.AgeRating = "R"
	}

	err = profileService.UpdateProfile(profile)
	if err != nil {
		fmt.Printf("\n[ERROR] %v\n", err)
	} else {
		fmt.Println("\n[OK] Perfil actualizado exitosamente")
	}
	time.Sleep(2 * time.Second)
}

func deleteProfile(profiles []models.Profile) {
	fmt.Print("\nIngrese el numero del perfil a eliminar (0 para cancelar): ")
	idx, err := utils.ToInt(utils.ReadLine(""))
	if err != nil || idx < 1 || idx > len(profiles) {
		return
	}

	profile := &profiles[idx-1]
	if profile.IsMain {
		fmt.Println("\n[ERROR] No se puede eliminar el perfil principal")
		time.Sleep(2 * time.Second)
		return
	}

	fmt.Printf("\nEsta seguro de eliminar el perfil '%s'? (S/N): ", profile.Name)
	confirm := strings.ToUpper(utils.ReadLine(""))
	if confirm != "S" {
		return
	}

	err = profileService.DeleteProfile(profile.ID)
	if err != nil {
		fmt.Printf("\n[ERROR] %v\n", err)
	} else {
		fmt.Println("\n[OK] Perfil eliminado exitosamente")
	}
	time.Sleep(2 * time.Second)
}

func getMaxProfilesForPlan(planID int) int {
	switch planID {
	case 1:
		return 1
	case 2:
		return 3
	case 3:
		return 5
	default:
		return 1
	}
}

func showMainMenu() {
	utils.ClearScreen()
	fmt.Println("=========================================")
	fmt.Printf("    Perfil: %s\n", currentProfile.Name)
	fmt.Println("=========================================")
	fmt.Println()
	fmt.Println("  [1] Inicio")
	fmt.Println("  [2] Tendencias")
	fmt.Println("  [3] Buscar Contenido")
	fmt.Println("  [4] Mi Lista")
	fmt.Println("  [5] Mis Playlists")
	fmt.Println("  [6] Historial")
	if currentUser.IsAdmin {
		fmt.Println("  [7] Panel de Administracion")
		fmt.Println("  [8] Cambiar Perfil")
		fmt.Println("  [0] Cerrar Sesion")
	} else {
		fmt.Println("  [7] Cambiar Perfil")
		fmt.Println("  [0] Cerrar Sesion")
	}
	fmt.Println()
	fmt.Print("Seleccione: ")

	option := utils.ReadLine("")
	switch option {
	case "1":
		showHome()
	case "2":
		showTrending()
	case "3":
		searchContent()
	case "4":
		showMyList()
	case "5":
		managePlaylists()
	case "6":
		viewHistory()
	case "7":
		if currentUser.IsAdmin {
			showAdminPanel()
		} else {
			currentProfile = nil
		}
	case "8":
		if currentUser.IsAdmin {
			currentProfile = nil
		}
	case "0":
		currentUser = nil
		currentProfile = nil
	default:
		fmt.Println("\n[ERROR] Opcion invalida")
		time.Sleep(1 * time.Second)
	}
}

func showHome() {
	utils.ClearScreen()
	fmt.Println("=========================================")
	fmt.Println("               Inicio                    ")
	fmt.Println("=========================================")
	fmt.Println()

	// Continuar viendo
	continueWatching, _ := playbackService.GetContinueWatching(currentProfile.ID)
	if len(continueWatching) > 0 {
		fmt.Println("--- Continuar viendo ---")
		for i, entry := range continueWatching {
			if i >= 5 {
				break
			}
			var title string
			if entry.ContentType == "audiovisual" {
				content, _ := contentService.GetAudiovisualByID(entry.ContentID)
				if content != nil {
					title = content.Title
					progress := (entry.Progress * 100) / (content.Duration * 60)
					fmt.Printf("  [%d] %s (%d%% visto)\n", entry.ContentID, title, progress)
				}
			} else {
				content, _ := contentService.GetAudioByID(entry.ContentID)
				if content != nil {
					title = content.Title
					progress := (entry.Progress * 100) / (content.Duration * 60)
					fmt.Printf("  [%d] %s - %s (%d%% escuchado)\n", entry.ContentID, content.Artist, title, progress)
				}
			}
		}
		fmt.Println()
	}

	// Recomendaciones
	fmt.Println("--- Recomendado para ti ---")
	recommendations := playbackService.GetRecommendations(currentProfile.ID, currentProfile.AgeRating)
	for i, content := range recommendations {
		if i >= 5 {
			break
		}
		switch c := content.(type) {
		case models.AudiovisualContent:
			fmt.Printf("  [%d] %s (%s) - %.1f/10\n", c.ID, c.Title, c.Type, c.AverageRating)
		case models.AudioContent:
			fmt.Printf("  [%d] %s - %s - %.1f/10\n", c.ID, c.Artist, c.Title, c.AverageRating)
		}
	}

	fmt.Println()
	fmt.Println("  [0] Volver")
	fmt.Print("\nSeleccione ID para reproducir (0 para volver): ")
	idStr := utils.ReadLine("")
	if idStr == "0" {
		return
	}
	
	id, err := utils.ToInt(idStr)
	if err == nil {
		playContentByID(id)
	}
}

func showTrending() {
	utils.ClearScreen()
	fmt.Println("=========================================")
	fmt.Println("             Tendencias                  ")
	fmt.Println("=========================================")
	fmt.Println()

	fmt.Println("--- Contenido Audiovisual Popular ---")
	audiovisuals, _ := contentService.GetAllAudiovisual()
	count := 0
	for _, av := range audiovisuals {
		if isContentAllowedForProfile(av.AgeRating, currentProfile.AgeRating) {
			fmt.Printf("  [%d] %s (%s) - %.1f/10\n", av.ID, av.Title, av.Genre, av.AverageRating)
			count++
			if count >= 5 {
				break
			}
		}
	}

	fmt.Println("\n--- Contenido de Audio Popular ---")
	audios, _ := contentService.GetAllAudio()
	count = 0
	for _, a := range audios {
		if isContentAllowedForProfile(a.AgeRating, currentProfile.AgeRating) {
			fmt.Printf("  [%d] %s - %s (%.1f/10)\n", a.ID, a.Artist, a.Title, a.AverageRating)
			count++
			if count >= 5 {
				break
			}
		}
	}

	fmt.Println()
	fmt.Println("  [0] Volver")
	fmt.Print("\nSeleccione ID para reproducir (0 para volver): ")
	idStr := utils.ReadLine("")
	if idStr == "0" {
		return
	}
	
	id, err := utils.ToInt(idStr)
	if err == nil {
		playContentByID(id)
	}
}

func searchContent() {
	utils.ClearScreen()
	fmt.Println("=========================================")
	fmt.Println("          Buscar Contenido               ")
	fmt.Println("=========================================")
	fmt.Println()
	fmt.Println("  [1] Buscar por Titulo")
	fmt.Println("  [2] Buscar por Genero")
	fmt.Println("  [3] Buscar por Actor")
	fmt.Println("  [0] Volver")
	fmt.Println()
	fmt.Print("Seleccione: ")

	option := utils.ReadLine("")
	switch option {
	case "1":
		searchByTitle()
	case "2":
		searchByGenre()
	case "3":
		searchByActor()
	case "0":
		return
	}
}

func searchByTitle() {
	utils.ClearScreen()
	fmt.Println("=========================================")
	fmt.Println("         Buscar por Titulo               ")
	fmt.Println("=========================================")
	fmt.Println()
	query := utils.ReadLine("Ingrese el titulo a buscar: ")
	if query == "" {
		return
	}

	results := contentService.SearchByTitle(query, currentProfile.AgeRating)
	if len(results) == 0 {
		fmt.Println("\n[INFO] No se encontraron resultados")
		time.Sleep(2 * time.Second)
		return
	}

	fmt.Println("\n--- Resultados ---")
	for _, content := range results {
		switch c := content.(type) {
		case models.AudiovisualContent:
			fmt.Printf("  [%d] %s (%s, %s) - %.1f/10\n", c.ID, c.Title, c.Type, c.Genre, c.AverageRating)
		case models.AudioContent:
			fmt.Printf("  [%d] %s - %s (%s) - %.1f/10\n", c.ID, c.Artist, c.Title, c.Type, c.AverageRating)
		}
	}

	fmt.Println()
	fmt.Print("Seleccione ID para ver detalles (0 para volver): ")
	idStr := utils.ReadLine("")
	if idStr == "0" {
		return
	}
	
	id, err := utils.ToInt(idStr)
	if err == nil {
		showContentDetails(id)
	}
}

func searchByGenre() {
	utils.ClearScreen()
	fmt.Println("=========================================")
	fmt.Println("         Buscar por Genero               ")
	fmt.Println("=========================================")
	fmt.Println()
	genre := utils.ReadLine("Ingrese el genero: ")
	if genre == "" {
		return
	}

	results := contentService.SearchByGenre(genre, currentProfile.AgeRating)
	if len(results) == 0 {
		fmt.Println("\n[INFO] No se encontraron resultados")
		time.Sleep(2 * time.Second)
		return
	}

	fmt.Println("\n--- Resultados ---")
	for _, content := range results {
		switch c := content.(type) {
		case models.AudiovisualContent:
			fmt.Printf("  [%d] %s (%s) - %.1f/10\n", c.ID, c.Title, c.Type, c.AverageRating)
		case models.AudioContent:
			fmt.Printf("  [%d] %s - %s - %.1f/10\n", c.ID, c.Artist, c.Title, c.AverageRating)
		}
	}

	fmt.Println()
	fmt.Print("Seleccione ID para ver detalles (0 para volver): ")
	idStr := utils.ReadLine("")
	if idStr == "0" {
		return
	}
	
	id, err := utils.ToInt(idStr)
	if err == nil {
		showContentDetails(id)
	}
}

func searchByActor() {
	utils.ClearScreen()
	fmt.Println("=========================================")
	fmt.Println("         Buscar por Actor                ")
	fmt.Println("=========================================")
	fmt.Println()
	actor := utils.ReadLine("Ingrese el nombre del actor: ")
	if actor == "" {
		return
	}

	results := contentService.SearchByActor(actor, currentProfile.AgeRating)
	if len(results) == 0 {
		fmt.Println("\n[INFO] No se encontraron resultados")
		time.Sleep(2 * time.Second)
		return
	}

	fmt.Println("\n--- Resultados ---")
	for _, av := range results {
		fmt.Printf("  [%d] %s (%s, %d) - %.1f/10\n", av.ID, av.Title, av.Genre, av.ReleaseYear, av.AverageRating)
	}

	fmt.Println()
	fmt.Print("Seleccione ID para ver detalles (0 para volver): ")
	idStr := utils.ReadLine("")
	if idStr == "0" {
		return
	}
	
	id, err := utils.ToInt(idStr)
	if err == nil {
		showContentDetails(id)
	}
}

func showContentDetails(contentID int) {
	// Intentar buscar en audiovisual
	av, err := contentService.GetAudiovisualByID(contentID)
	if err == nil && av != nil {
		showAudiovisualDetails(av)
		return
	}

	// Intentar buscar en audio
	audio, err := contentService.GetAudioByID(contentID)
	if err == nil && audio != nil {
		showAudioDetails(audio)
		return
	}

	fmt.Println("\n[ERROR] Contenido no encontrado")
	time.Sleep(2 * time.Second)
}

func showAudiovisualDetails(content *models.AudiovisualContent) {
	for {
		utils.ClearScreen()
		fmt.Println("=========================================")
		fmt.Printf("        %s\n", content.Title)
		fmt.Println("=========================================")
		fmt.Println()
		fmt.Printf("Tipo: %s\n", content.Type)
		fmt.Printf("Genero: %s\n", content.Genre)
		fmt.Printf("Director: %s\n", content.Director)
		if content.Actors != "" {
			fmt.Printf("Actores: %s\n", content.Actors)
		}
		fmt.Printf("Ano: %d\n", content.ReleaseYear)
		fmt.Printf("Duracion: %d minutos\n", content.Duration)
		fmt.Printf("Clasificacion: %s\n", content.AgeRating)
		fmt.Printf("Calificacion: %.1f/10\n", content.AverageRating)
		fmt.Printf("\nSinopsis:\n%s\n", content.Synopsis)
		fmt.Println()
		fmt.Println("  [1] Reproducir")
		fmt.Println("  [2] Agregar a Mi Lista")
		fmt.Println("  [3] Agregar a Playlist")
		fmt.Println("  [4] Calificar")
		fmt.Println("  [0] Volver")
		fmt.Println()
		fmt.Print("Seleccione: ")

		option := utils.ReadLine("")
		switch option {
		case "1":
			playAudiovisual(content.ID)
			return
		case "2":
			err := playbackService.AddFavorite(currentProfile.ID, content.ID, "audiovisual")
			if err != nil {
				fmt.Printf("\n[ERROR] %v\n", err)
			} else {
				fmt.Println("\n[OK] Agregado a Mi Lista")
			}
			time.Sleep(2 * time.Second)
		case "3":
			addToPlaylist(content.ID, "audiovisual")
		case "4":
			rateContent(content.ID, "audiovisual")
		case "0":
			return
		}
	}
}

func showAudioDetails(content *models.AudioContent) {
	for {
		utils.ClearScreen()
		fmt.Println("=========================================")
		fmt.Printf("        %s\n", content.Title)
		fmt.Println("=========================================")
		fmt.Println()
		fmt.Printf("Artista: %s\n", content.Artist)
		fmt.Printf("Album: %s\n", content.Album)
		fmt.Printf("Tipo: %s\n", content.Type)
		fmt.Printf("Genero: %s\n", content.Genre)
		fmt.Printf("Duracion: %d minutos\n", content.Duration)
		fmt.Printf("Clasificacion: %s\n", content.AgeRating)
		fmt.Printf("Calificacion: %.1f/10\n", content.AverageRating)
		fmt.Println()
		fmt.Println("  [1] Reproducir")
		fmt.Println("  [2] Agregar a Mi Lista")
		fmt.Println("  [3] Agregar a Playlist")
		fmt.Println("  [4] Calificar")
		fmt.Println("  [0] Volver")
		fmt.Println()
		fmt.Print("Seleccione: ")

		option := utils.ReadLine("")
		switch option {
		case "1":
			playAudio(content.ID)
			return
		case "2":
			err := playbackService.AddFavorite(currentProfile.ID, content.ID, "audio")
			if err != nil {
				fmt.Printf("\n[ERROR] %v\n", err)
			} else {
				fmt.Println("\n[OK] Agregado a Mi Lista")
			}
			time.Sleep(2 * time.Second)
		case "3":
			addToPlaylist(content.ID, "audio")
		case "4":
			rateContent(content.ID, "audio")
		case "0":
			return
		}
	}
}

func showAudioDetails(content *models.AudioContent) {
	for {
		utils.ClearScreen()
		fmt.Println("=========================================")
		fmt.Printf("        %s\n", content.Title)
		fmt.Println("=========================================")
		fmt.Println()
		fmt.Printf("Artista: %s\n", content.Artist)
		fmt.Printf("Album: %s\n", content.Album)
		fmt.Printf("Tipo: %s\n", content.Type)
		fmt.Printf("Genero: %s\n", content.Genre)
		fmt.Printf("Duracion: %d minutos\n", content.Duration)
		fmt.Printf("Clasificacion: %s\n", content.AgeRating)
		fmt.Printf("Calificacion: %.1f/10\n", content.AverageRating)
		fmt.Println()
		fmt.Println("  [1] Reproducir")
		fmt.Println("  [2] Agregar a Mi Lista")
		fmt.Println("  [3] Agregar a Playlist")
		fmt.Println("  [4] Calificar")
		fmt.Println("  [0] Volver")
		fmt.Println()
		fmt.Print("Seleccione: ")

		option := utils.ReadLine("")
		switch option {
		case "1":
			playAudio(content.ID)
			return
		case "2":
			err := playbackService.AddFavorite(currentProfile.ID, content.ID, "audio")
			if err != nil {
				fmt.Printf("\n[ERROR] %v\n", err)
			} else {
				fmt.Println("\n[OK] Agregado a Mi Lista")
			}
			time.Sleep(2 * time.Second)
		case "3":
			addToPlaylist(content.ID, "audio")
		case "4":
			rateContent(content.ID, "audio")
		case "0":
			return
		}
	}
}

func playContentByID(contentID int) {
	av, err := contentService.GetAudiovisualByID(contentID)
	if err == nil && av != nil {
		playAudiovisual(contentID)
		return
	}

	audio, err := contentService.GetAudioByID(contentID)
	if err == nil && audio != nil {
		playAudio(contentID)
		return
	}

	fmt.Println("\n[ERROR] Contenido no encontrado")
	time.Sleep(2 * time.Second)
}

func playAudiovisual(contentID int) {
	content, err := contentService.GetAudiovisualByID(contentID)
	if err != nil {
		fmt.Println("\n[ERROR] Error al cargar contenido")
		time.Sleep(2 * time.Second)
		return
	}

	utils.ClearScreen()
	fmt.Println("=========================================")
	fmt.Printf("        Reproduciendo: %s\n", content.Title)
	fmt.Println("=========================================")
	fmt.Println()
	fmt.Println("Simulando reproduccion...")
	fmt.Println("[####################] 100%")
	fmt.Printf("Duracion total: %d minutos\n", content.Duration)
	fmt.Println("=========================================")

	playbackService.AddToHistory(currentProfile.ID, contentID, "audiovisual")
	progressSeconds := (content.Duration * 60) / 2
	playbackService.UpdateProgress(currentProfile.ID, contentID, "audiovisual", progressSeconds)

	fmt.Println("\n[OK] Reproduccion finalizada")
	fmt.Println("Se ha guardado tu progreso.")
	time.Sleep(3 * time.Second)
}

func playAudio(contentID int) {
	content, err := contentService.GetAudioByID(contentID)
	if err != nil {
		fmt.Println("\n[ERROR] Error al cargar contenido")
		time.Sleep(2 * time.Second)
		return
	}

	utils.ClearScreen()
	fmt.Println("=========================================")
	fmt.Printf("        Reproduciendo: %s\n", content.Title)
	fmt.Println("=========================================")
	fmt.Println()
	fmt.Printf("Artista: %s\n", content.Artist)
	fmt.Println("\nSimulando reproduccion...")
	fmt.Println("[####################] 100%")
	fmt.Printf("Duracion total: %d minutos\n", content.Duration)
	fmt.Println("=========================================")

	playbackService.AddToHistory(currentProfile.ID, contentID, "audio")
	progressSeconds := (content.Duration * 60) * 7 / 10
	playbackService.UpdateProgress(currentProfile.ID, contentID, "audio", progressSeconds)

	fmt.Println("\n[OK] Reproduccion finalizada")
	fmt.Println("Se ha guardado tu progreso.")
	time.Sleep(3 * time.Second)
}

func rateContent(contentID int, contentType string) {
	utils.ClearScreen()
	fmt.Println("=========================================")
	fmt.Println("         Calificar Contenido             ")
	fmt.Println("=========================================")
	fmt.Println()
	fmt.Println("Ingrese su calificacion (1.0 - 10.0)")
	ratingStr := utils.ReadLine("Calificacion: ")

	rating, err := utils.ToFloat(ratingStr)
	if err != nil || rating < 1.0 || rating > 10.0 {
		fmt.Println("\n[ERROR] Calificacion invalida. Debe ser entre 1.0 y 10.0")
		time.Sleep(2 * time.Second)
		return
	}

	err = contentService.RateContent(currentProfile.ID, contentID, contentType, rating)
	if err != nil {
		fmt.Printf("\n[ERROR] %v\n", err)
	} else {
		fmt.Printf("\n[OK] Has calificado este contenido con %.1f/10\n", rating)
	}
	time.Sleep(2 * time.Second)
}

func showMyList() {
	utils.ClearScreen()
	fmt.Println("=========================================")
	fmt.Println("              Mi Lista                   ")
	fmt.Println("=========================================")
	fmt.Println()

	favorites, err := playbackService.GetFavorites(currentProfile.ID)
	if err != nil {
		fmt.Printf("[ERROR] %v\n", err)
		time.Sleep(2 * time.Second)
		return
	}

	if len(favorites) == 0 {
		fmt.Println("[INFO] No tienes ningun contenido en tu lista")
		fmt.Println()
		fmt.Println("  [0] Volver")
		utils.ReadLine("")
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
				fmt.Printf("  [%d] %s\n", content.ID, title)
				fmt.Printf("       %s\n", details)
			}
		} else {
			content, _ := contentService.GetAudioByID(fav.ContentID)
			if content != nil {
				title = fmt.Sprintf("%s - %s", content.Artist, content.Title)
				details = fmt.Sprintf("[%s] %s", content.Type, content.Genre)
				fmt.Printf("  [%d] %s\n", content.ID, title)
				fmt.Printf("       %s\n", details)
			}
		}
	}
	fmt.Println()
	fmt.Println("  [0] Volver")
	fmt.Print("\nSeleccione ID para reproducir (0 para volver): ")
	idStr := utils.ReadLine("")
	if idStr == "0" {
		return
	}
	
	id, err := utils.ToInt(idStr)
	if err == nil {
		playContentByID(id)
	}
}

func managePlaylists() {
	for {
		utils.ClearScreen()
		fmt.Println("=========================================")
		fmt.Println("           Mis Playlists                 ")
		fmt.Println("=========================================")
		fmt.Println()

		playlists, _ := playlistService.GetPlaylistsByProfileID(currentProfile.ID)
		
		if len(playlists) == 0 {
			fmt.Println("[INFO] No tienes playlists creadas")
		} else {
			for i, p := range playlists {
				fmt.Printf("  [%d] %s\n", i+1, p.Name)
				if p.Description != "" {
					fmt.Printf("       %s\n", p.Description)
				}
			}
		}

		fmt.Println()
		fmt.Println("  [A] Crear Playlist")
		fmt.Println("  [V] Ver Playlist")
		fmt.Println("  [D] Eliminar Playlist")
		fmt.Println("  [0] Volver")
		fmt.Println()
		fmt.Print("Seleccione: ")

		option := strings.ToUpper(utils.ReadLine(""))
		switch option {
		case "A":
			createPlaylist()
		case "V":
			viewPlaylist(playlists)
		case "D":
			deletePlaylist(playlists)
		case "0":
			return
		}
	}
}

func createPlaylist() {
	utils.ClearScreen()
	fmt.Println("=========================================")
	fmt.Println("          Crear Playlist                 ")
	fmt.Println("=========================================")
	fmt.Println()

	name := utils.ReadLine("Nombre de la playlist: ")
	if name == "" {
		fmt.Println("\n[ERROR] El nombre no puede estar vacio")
		time.Sleep(2 * time.Second)
		return
	}

	description := utils.ReadLine("Descripcion (opcional): ")

	_, err := playlistService.CreatePlaylist(currentProfile.ID, name, description)
	if err != nil {
		fmt.Printf("\n[ERROR] %v\n", err)
	} else {
		fmt.Println("\n[OK] Playlist creada exitosamente")
	}
	time.Sleep(2 * time.Second)
}

func viewPlaylist(playlists []models.Playlist) {
	if len(playlists) == 0 {
		return
	}

	fmt.Print("\nIngrese el numero de la playlist (0 para cancelar): ")
	idx, err := utils.ToInt(utils.ReadLine(""))
	if err != nil || idx < 1 || idx > len(playlists) {
		return
	}

	playlist := playlists[idx-1]
	for {
		utils.ClearScreen()
		fmt.Println("=========================================")
		fmt.Printf("        Playlist: %s\n", playlist.Name)
		fmt.Println("=========================================")
		fmt.Println()
		if playlist.Description != "" {
			fmt.Printf("Descripcion: %s\n\n", playlist.Description)
		}

		items, _ := playlistService.GetPlaylistItems(playlist.ID)
		if len(items) == 0 {
			fmt.Println("[INFO] Esta playlist esta vacia")
		} else {
			fmt.Println("Contenido:")
			for _, item := range items {
				if item.ContentType == "audiovisual" {
					content, _ := contentService.GetAudiovisualByID(item.ContentID)
					if content != nil {
						fmt.Printf("  [%d] %s (%s)\n", content.ID, content.Title, content.Type)
					}
				} else {
					content, _ := contentService.GetAudioByID(item.ContentID)
					if content != nil {
						fmt.Printf("  [%d] %s - %s\n", content.ID, content.Artist, content.Title)
					}
				}
			}
		}

		fmt.Println()
		fmt.Println("  [P] Reproducir contenido")
		fmt.Println("  [R] Eliminar contenido")
		fmt.Println("  [0] Volver")
		fmt.Println()
		fmt.Print("Seleccione: ")

		option := strings.ToUpper(utils.ReadLine(""))
		switch option {
		case "P":
			fmt.Print("\nIngrese ID del contenido: ")
			idStr := utils.ReadLine("")
			id, err := utils.ToInt(idStr)
			if err == nil {
				playContentByID(id)
			}
		case "R":
			fmt.Print("\nIngrese ID del contenido a eliminar: ")
			idStr := utils.ReadLine("")
			id, err := utils.ToInt(idStr)
			if err == nil {
				// Buscar tipo de contenido
				contentType := "audiovisual"
				_, err := contentService.GetAudiovisualByID(id)
				if err != nil {
					contentType = "audio"
				}
				err = playlistService.RemoveItemFromPlaylist(playlist.ID, id, contentType)
				if err != nil {
					fmt.Printf("\n[ERROR] %v\n", err)
				} else {
					fmt.Println("\n[OK] Contenido eliminado de la playlist")
				}
				time.Sleep(2 * time.Second)
			}
		case "0":
			return
		}
	}
}

func deletePlaylist(playlists []models.Playlist) {
	if len(playlists) == 0 {
		return
	}

	fmt.Print("\nIngrese el numero de la playlist a eliminar (0 para cancelar): ")
	idx, err := utils.ToInt(utils.ReadLine(""))
	if err != nil || idx < 1 || idx > len(playlists) {
		return
	}

	playlist := playlists[idx-1]
	fmt.Printf("\nEsta seguro de eliminar '%s'? (S/N): ", playlist.Name)
	confirm := strings.ToUpper(utils.ReadLine(""))
	if confirm != "S" {
		return
	}

	err = playlistService.DeletePlaylist(playlist.ID)
	if err != nil {
		fmt.Printf("\n[ERROR] %v\n", err)
	} else {
		fmt.Println("\n[OK] Playlist eliminada exitosamente")
	}
	time.Sleep(2 * time.Second)
}

func addToPlaylist(contentID int, contentType string) {
	playlists, _ := playlistService.GetPlaylistsByProfileID(currentProfile.ID)
	if len(playlists) == 0 {
		fmt.Println("\n[INFO] No tienes playlists. Crea una primero.")
		time.Sleep(2 * time.Second)
		return
	}

	utils.ClearScreen()
	fmt.Println("=========================================")
	fmt.Println("       Agregar a Playlist                ")
	fmt.Println("=========================================")
	fmt.Println()
	for i, p := range playlists {
		fmt.Printf("  [%d] %s\n", i+1, p.Name)
	}
	fmt.Println()
	fmt.Print("Seleccione numero de playlist (0 para cancelar): ")
	idx, err := utils.ToInt(utils.ReadLine(""))
	if err != nil || idx < 1 || idx > len(playlists) {
		return
	}

	playlist := playlists[idx-1]
	err = playlistService.AddItemToPlaylist(playlist.ID, contentID, contentType)
	if err != nil {
		fmt.Printf("\n[ERROR] %v\n", err)
	} else {
		fmt.Printf("\n[OK] Agregado a '%s'\n", playlist.Name)
	}
	time.Sleep(2 * time.Second)
}

func viewHistory() {
	utils.ClearScreen()
	fmt.Println("=========================================")
	fmt.Println("       Historial de Reproduccion         ")
	fmt.Println("=========================================")
	fmt.Println()

	history, err := playbackService.GetHistory(currentProfile.ID)
	if err != nil {
		fmt.Printf("[ERROR] %v\n", err)
		time.Sleep(2 * time.Second)
		return
	}

	if len(history) == 0 {
		fmt.Println("[INFO] No tienes historial de reproduccion")
		fmt.Println()
		fmt.Println("  [0] Volver")
		utils.ReadLine("")
		return
	}

	fmt.Println("Tus ultimas reproducciones:")
	for i, entry := range history {
		if i >= 20 {
			break
		}
		var title string
		if entry.ContentType == "audiovisual" {
			content, _ := contentService.GetAudiovisualByID(entry.ContentID)
			if content != nil {
				title = content.Title
				progress := 0
				if content.Duration > 0 {
					progress = (entry.Progress * 100) / (content.Duration * 60)
				}
				fmt.Printf("  [%d] %s (%d%% visto)\n", content.ID, title, progress)
			}
		} else {
			content, _ := contentService.GetAudioByID(entry.ContentID)
			if content != nil {
				title = content.Title
				progress := 0
				if content.Duration > 0 {
					progress = (entry.Progress * 100) / (content.Duration * 60)
				}
				fmt.Printf("  [%d] %s - %s (%d%%)\n", content.ID, content.Artist, title, progress)
			}
		}
	}
	fmt.Println()
	fmt.Println("  [0] Volver")
	fmt.Print("\nSeleccione ID para reproducir (0 para volver): ")
	idStr := utils.ReadLine("")
	if idStr == "0" {
		return
	}
	
	id, err := utils.ToInt(idStr)
	if err == nil {
		playContentByID(id)
	}
}

func isContentAllowedForProfile(contentRating, profileRating string) bool {
	ratings := map[string]int{"G": 1, "PG": 2, "PG-13": 3, "R": 4, "General": 1, "Explicit": 4}
	contentLevel := ratings[contentRating]
	profileLevel := ratings[profileRating]
	return contentLevel <= profileLevel
}

// Continuacion del main.go - Panel de Administracion

func showAdminPanel() {
	for {
		utils.ClearScreen()
		fmt.Println("=========================================")
		fmt.Println("      Panel de Administracion            ")
		fmt.Println("=========================================")
		fmt.Println()
		fmt.Println("  [1] Gestionar Usuarios")
		fmt.Println("  [2] Gestionar Contenido")
		fmt.Println("  [3] Generar Reportes")
		fmt.Println("  [0] Volver")
		fmt.Println()
		fmt.Print("Seleccione: ")

		option := utils.ReadLine("")
		switch option {
		case "1":
			manageUsers()
		case "2":
			manageContent()
		case "3":
			generateReports()
		case "0":
			return
		default:
			fmt.Println("\n[ERROR] Opcion invalida")
			time.Sleep(1 * time.Second)
		}
	}
}

func manageUsers() {
	utils.ClearScreen()
	fmt.Println("=========================================")
	fmt.Println("         Gestion de Usuarios             ")
	fmt.Println("=========================================")
	fmt.Println()
	
	users, err := userService.GetAllUsers()
	if err != nil {
		fmt.Printf("[ERROR] %v\n", err)
		time.Sleep(2 * time.Second)
		return
	}

	for _, u := range users {
		planName := getPlanName(u.PlanID)
		adminTag := ""
		if u.IsAdmin {
			adminTag = " [ADMIN]"
		}
		fmt.Printf("- %s%s (%s)\n", u.Name, adminTag, u.Email)
		fmt.Printf("  Plan: %s | Edad: %d | Registro: %s\n", planName, u.Age, u.CreatedAt.Format("2006-01-02"))
		fmt.Println()
	}

	fmt.Println("  [0] Volver")
	utils.ReadLine("")
}

func manageContent() {
	for {
		utils.ClearScreen()
		fmt.Println("=========================================")
		fmt.Println("        Gestion de Contenido             ")
		fmt.Println("=========================================")
		fmt.Println()
		fmt.Println("  [1] Agregar Contenido Audiovisual")
		fmt.Println("  [2] Agregar Contenido de Audio")
		fmt.Println("  [3] Listar Contenido Audiovisual")
		fmt.Println("  [4] Listar Contenido de Audio")
		fmt.Println("  [0] Volver")
		fmt.Println()
		fmt.Print("Seleccione: ")

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
		case "0":
			return
		default:
			fmt.Println("\n[ERROR] Opcion invalida")
			time.Sleep(1 * time.Second)
		}
	}
}

func addAudiovisualContent() {
	utils.ClearScreen()
	fmt.Println("=========================================")
	fmt.Println("   Agregar Contenido Audiovisual         ")
	fmt.Println("=========================================")
	fmt.Println()

	title := utils.ReadLine("Titulo: ")
	contentType := utils.ReadLine("Tipo (movie/series/documentary): ")
	genre := utils.ReadLine("Genero: ")
	durationStr := utils.ReadLine("Duracion (minutos): ")
	duration, err := utils.ToInt(durationStr)
	if err != nil {
		fmt.Println("\n[ERROR] Duracion invalida")
		time.Sleep(2 * time.Second)
		return
	}

	ageRating := utils.ReadLine("Clasificacion (G/PG/PG-13/R): ")
	synopsis := utils.ReadLine("Sinopsis: ")
	yearStr := utils.ReadLine("Ano de lanzamiento: ")
	year, err := utils.ToInt(yearStr)
	if err != nil {
		fmt.Println("\n[ERROR] Ano invalido")
		time.Sleep(2 * time.Second)
		return
	}
	director := utils.ReadLine("Director: ")
	actors := utils.ReadLine("Actores (separados por coma): ")

	err = contentService.CreateAudiovisual(title, contentType, genre, duration, ageRating, synopsis, year, director, actors)
	if err != nil {
		fmt.Printf("\n[ERROR] %v\n", err)
	} else {
		fmt.Println("\n[OK] Contenido agregado exitosamente")
	}
	time.Sleep(2 * time.Second)
}

func addAudioContent() {
	utils.ClearScreen()
	fmt.Println("=========================================")
	fmt.Println("      Agregar Contenido de Audio         ")
	fmt.Println("=========================================")
	fmt.Println()

	title := utils.ReadLine("Titulo: ")
	contentType := utils.ReadLine("Tipo (song/podcast/audiobook): ")
	genre := utils.ReadLine("Genero: ")
	durationStr := utils.ReadLine("Duracion (minutos): ")
	duration, err := utils.ToInt(durationStr)
	if err != nil {
		fmt.Println("\n[ERROR] Duracion invalida")
		time.Sleep(2 * time.Second)
		return
	}

	ageRating := utils.ReadLine("Clasificacion (General/Explicit): ")
	artist := utils.ReadLine("Artista: ")
	album := utils.ReadLine("Album: ")
	trackStr := utils.ReadLine("Numero de pista: ")
	trackNumber, err := utils.ToInt(trackStr)
	if err != nil {
		trackNumber = 1
	}

	err = contentService.CreateAudio(title, contentType, genre, duration, ageRating, artist, album, trackNumber)
	if err != nil {
		fmt.Printf("\n[ERROR] %v\n", err)
	} else {
		fmt.Println("\n[OK] Contenido agregado exitosamente")
	}
	time.Sleep(2 * time.Second)
}

func listAudiovisualAdmin() {
	utils.ClearScreen()
	fmt.Println("=========================================")
	fmt.Println("   Lista de Contenido Audiovisual        ")
	fmt.Println("=========================================")
	fmt.Println()

	contents, err := contentService.GetAllAudiovisual()
	if err != nil {
		fmt.Printf("[ERROR] %v\n", err)
		time.Sleep(2 * time.Second)
		return
	}

	for _, c := range contents {
		fmt.Printf("[%d] %s (%s)\n", c.ID, c.Title, c.Type)
		fmt.Printf("     Genero: %s | Duracion: %d min | Rating: %.1f/10\n", c.Genre, c.Duration, c.AverageRating)
		fmt.Println()
	}

	fmt.Println("  [0] Volver")
	utils.ReadLine("")
}

func listAudioAdmin() {
	utils.ClearScreen()
	fmt.Println("=========================================")
	fmt.Println("      Lista de Contenido de Audio        ")
	fmt.Println("=========================================")
	fmt.Println()

	contents, err := contentService.GetAllAudio()
	if err != nil {
		fmt.Printf("[ERROR] %v\n", err)
		time.Sleep(2 * time.Second)
		return
	}

	for _, c := range contents {
		fmt.Printf("[%d] %s - %s (%s)\n", c.ID, c.Artist, c.Title, c.Type)
		fmt.Printf("     Album: %s | Duracion: %d min | Rating: %.1f/10\n", c.Album, c.Duration, c.AverageRating)
		fmt.Println()
	}

	fmt.Println("  [0] Volver")
	utils.ReadLine("")
}

func generateReports() {
	for {
		utils.ClearScreen()
		fmt.Println("=========================================")
		fmt.Println("          Generar Reportes               ")
		fmt.Println("=========================================")
		fmt.Println()
		fmt.Println("  [1] Reporte de Usuarios")
		fmt.Println("  [2] Reporte de Contenido")
		fmt.Println("  [3] Reporte de Ingresos")
		fmt.Println("  [0] Volver")
		fmt.Println()
		fmt.Print("Seleccione: ")

		option := utils.ReadLine("")
		switch option {
		case "1":
			showUserReport()
		case "2":
			showContentReport()
		case "3":
			showRevenueReport()
		case "0":
			return
		default:
			fmt.Println("\n[ERROR] Opcion invalida")
			time.Sleep(1 * time.Second)
		}
	}
}

func showUserReport() {
	utils.ClearScreen()
	fmt.Println("=========================================")
	fmt.Println("         Reporte de Usuarios             ")
	fmt.Println("=========================================")
	fmt.Println()

	report, err := reportService.GenerateUserReport()
	if err != nil {
		fmt.Printf("[ERROR] %v\n", err)
		time.Sleep(2 * time.Second)
		return
	}

	fmt.Printf("Total de Usuarios: %d\n\n", report.TotalUsers)
	fmt.Println("Usuarios por Plan:")
	for plan, count := range report.UsersByPlan {
		fmt.Printf("  - %s: %d usuarios\n", plan, count)
	}

	fmt.Println()
	fmt.Println("  [0] Volver")
	utils.ReadLine("")
}

func showContentReport() {
	utils.ClearScreen()
	fmt.Println("=========================================")
	fmt.Println("         Reporte de Contenido            ")
	fmt.Println("=========================================")
	fmt.Println()

	report, err := reportService.GenerateContentReport()
	if err != nil {
		fmt.Printf("[ERROR] %v\n", err)
		time.Sleep(2 * time.Second)
		return
	}

	fmt.Printf("Total Audiovisual: %d\n", report.TotalAudiovisual)
	fmt.Printf("Total Audio: %d\n\n", report.TotalAudio)

	if len(report.TopRatedContent) > 0 {
		fmt.Println("Contenido Mejor Calificado:")
		for _, item := range report.TopRatedContent {
			fmt.Printf("  - %s\n", item)
		}
	}

	fmt.Println()
	fmt.Println("  [0] Volver")
	utils.ReadLine("")
}

func showRevenueReport() {
	utils.ClearScreen()
	fmt.Println("=========================================")
	fmt.Println("         Reporte de Ingresos             ")
	fmt.Println("=========================================")
	fmt.Println()

	report, err := reportService.GenerateRevenueReport()
	if err != nil {
		fmt.Printf("[ERROR] %v\n", err)
		time.Sleep(2 * time.Second)
		return
	}

	fmt.Printf("Ingresos Totales: $%.2f\n", report.TotalRevenue)
	fmt.Printf("Total Transacciones: %d\n\n", report.Transactions)

	if len(report.RevenueByPlan) > 0 {
		fmt.Println("Ingresos por Plan:")
		for plan, revenue := range report.RevenueByPlan {
			fmt.Printf("  - %s: $%.2f\n", plan, revenue)
		}
	}

	fmt.Println()
	fmt.Println("  [0] Volver")
	utils.ReadLine("")
}

func getPlanName(planID int) string {
	switch planID {
	case 1:
		return "Free"
	case 2:
		return "Estandar"
	case 3:
		return "Premium 4K"
	default:
		return "Desconocido"
	}
}