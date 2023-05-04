package main

type menu struct{}

func (menu) roomAction() string {
	return "1. Создать комнату\n" +
		"2. Присоедениться к комнате"
}

func (menu) createRoomTitle() string {
	return "Введите название"
}

func (menu) createRoomPass() string {
	return "Введите пароль"
}

func (menu) roomCreated() string {
	return "Комната создана"
}

func (menu) joinRoomTitle() string {
	return "Введите название"
}

func (menu) joinRoomPass() string {
	return "Введите пароль"
}

func (menu) joinRoomComp() string {
	return "Вы успешно вошли в комнату"
}

func (menu) noSuchAnswer() string {
	return "Нет такого варианта ответа"
}
