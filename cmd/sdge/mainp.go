package mainp
func showMenu(title string, options map[string]string, handler func(string)) {
    for {
        utils.ClearScreen()
        fmt.Println("=========================================")
        fmt.Println(" ", title)
        fmt.Println("=========================================")
        fmt.Println()

        for key, desc := range options {
            fmt.Printf("  [%s] %s\n", key, desc)
        }
        fmt.Println()
        fmt.Print("Seleccione: ")

        choice := strings.ToUpper(utils.ReadLine(""))
        handler(choice)

        if choice == "0" {
            break
        }
    }
}
func showAuthMenu() {
    options := map[string]string{
        "1": "Iniciar Sesion",
        "2": "Registrarse",
        "0": "Salir",
    }

    showMenu("SDGEStreaming - Bienvenido", options, func(choice string) {
        switch choice {
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
    })
}
func ReadOption(prompt string, validOptions []string) string {
    for {
        fmt.Print(prompt)
        input := strings.ToUpper(utils.ReadLine(""))
        for _, v := range validOptions {
            if input == v {
                return input
            }
        }
        fmt.Println("[ERROR] Opcion invalida")
    }
}
options := map[string]string{
    "1": "Inicio",
    "2": "Tendencias",
    "3": "Buscar Contenido",
    "0": "Cerrar Sesion",
}
showMenu("Menu Principal", options, handleMainMenuOption)
