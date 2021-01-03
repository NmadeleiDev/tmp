зареганный домен (сейчас опущен):
https://aim-love.ga

КАК ПОДНЯТЬ

1. Скопировать .env.example в .env файл, заполнить поля SERVICE_MAIL_ADDR и SERVICE_MAIL_PASSWD (если хочешь, чтобы приходили письма с подтверждением). Там уже есть настроенный gmail для примера, пароль можешь спросить у меня.
2. Написать "make up" в корне. Если захочешь поменять порты - все в .env. Сейчас апи поднимается на порте 8080. В принципе, после make up можешь пойти попить чай или что такое - так как будет подниматься сервер и две базы данных (это не быстро, но только в первый раз, потом докер сохранит образы и будет быстро). Возможно, там постгрес с первого раза не успеет подняться, так что если увидешь, что не работает - сделай make down && make up.


КАК ПОЛЬЗОВАТЬСЯ

В данный реализовано:

- POST /api/main/signup - после этого запроса, на почту, указанную в email придет письмо с ссылкой для подтверждения (с почты, указанной в .env SERVICE_MAIL_ADDR)
Валидное тело запроса:
{
	"email": "test3@gmail.com",
	"phone": "89671102000",
	"password": "123",
	"username": "Mary",
	"birthDate": 958003200,
	"gender": "female",
	"country": "Russia",
	"city": "Moscow",
	"maxDist": 100,
	"lookFor": "male",
	"minAge": 17,
	"maxAge": 38,
}

- POST /api/main/signin - в ответ будет вся инфа юзера!
Валидное тело:
{
	"email": "hello@gmail.com",
	"password": "123"
}


ВОССТАНОВЛЕНИЕ ПАРОЛЯ
- POST /api/main/reset - инициализация процедуры восстановления. При этом запросе на сервере создаются необходимые ключи и на почту пользователю отправляется письмо для сброса пароля
Тело:
{
    "email": "<меил юзера>"
}

- GET /api/main/reset - проверка валидности попытки (отправлять сразу при переходе на страницу для смены пароля). При успешной проверке в ответ ставится httpOnly кука для аутентификации самой смены пароля
Параметры: k - ключ, взятый из одноименного параметра изначальной ссылки

- PUT /api/main/reset - смена пароля. Успешной будет только при наличии валидной куки.
Тело:
{
    "password": "новый пароль"
}



- POST /api/main/user - обновление данных пользователя. Тут просто все данные пользователя обновятся на то, что отправлено. То есть нужно все эти поля обязательно присылать, если прислать какие - то пустые - они станут пустыми. Я думаю, так удобнее, так как на фронте все равно будут уже после signin все актуальные данные, и сюда ты просто их же присылаешь, изменяя те, которые поменял юзер. Я имею ввиду, просто реактивно выводишь на страничку, реактивно они обновляются у тебя, и если юзер что то сохраняет - отсылаешь такой запрос на обновление.ы 
Тело запроса:
{
	"id": "user id",
	"phone": "89671102000",
	"username": "Liza",
	"name": "Liza",
	"surname": "Liza",
	"age": 17,
	"gender": "female",
	"country": "Russia",
	"city": "Moscow",
	"maxDist": 100,
	"lookFor": "male",
	"minAge": 24,
	"maxAge": 47
}

- GET /api/main/account - получение собственных данных пользователя

- DELETE /api/main/signout - тут ничего, кроме самого запроса отправлять не нужно, сервер просто сам обнулит сессию

- DELETE /api/main/account - удаление аккаунта

- GET /api/main/strangers (если придумаешь название получше - супер). Это основной запрос, выполняющийся после того, как человек залогинился. Он возвращает ему пачку пользователей сайта, подходяший по его критериям (возраст, город, пол). Типо лента в сайте знакомств.

- GET /api/main/data/{id юзера} - получение публичных данных юзера по его id. GET параметр full=false позволяет получить короткую версию данных. Во всех иных случаях высылаются все публичные данные


Бан запросы /api/main/ban:
- GET - в ответ придет массив id пользователей, забанненых юзером
- POST - тело: {id: <id аккаунта для забанивания>}
- DELETE - тело: {id: <id аккаунта для разбанивания>}


MEDIA запросы

- POST /api/media/upload - загрузка фоток. Тело:
{
    isAvatar: Boolean (true можешь писать только, когда устанавливаешь аватар, при любом другом значении этого поля, даже при его отсутсвии, фотка просто сохраниться в галерею пользователя)
    userImage: <файл картинки (для того чтобы это был файл, надо просто указать тип инпута file)>
}

- PUT /api/media/avatar - установка уже загруженной фотки как аватара. Тело:
{
    imageId: <id существующей картинки этого пользователя>
}

- GET /api/media/img/<id картинки> - получение картинки

- DELETE /api/media/img - удаление фоток. Тело:
{
    images: <id картинок к удалению>
}


ГРУППА ЗАПРОСОВ ДЛЯ СОЦ ИНТЕРАКЦИЙ:

/api/main/look
- POST - в теле отправлять id того, кого посмотрел
- GET - получить список тех, кого посмотрел

/api/main/like
- POST - в теле отправлять id того, кого лайкнул
- GET - получить список тех, кого лайкнул
- DELETE (+ get параметр id=<id к удалению>) - убрать лайк с чела. Если там был матч, он тоже уберется


Пример тела, содежащего id:
{
	"id": "kjhaskfhs87yf9ay94hkhf2k298e"
}



Запросы отправлять на http://localhost:8080