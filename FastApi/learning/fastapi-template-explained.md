# FastAPI Template Loyihasi — To'liq Tahlil

Bu hujjat `fastapi-structure` loyihasini **fayl-fayl** tahlil qiladi. Loyiha real "Foreign Trip / Tashqi Xizmat Safari" (UzAssets) tizimidan olingan production-grade FastAPI template hisoblanadi.

> Maqsad: tashkilotlar tashqi xizmat safariga chiqish uchun tizimga kirib loyiha yaratadi, har bir loyiha MinFin Ijro Apparati Departamentiga keladi, departament xodimlari kelishadi, oxirida G'aznachilik Qo'mitasi tasdiqlaydi.

## Texnologiyalar to'plami

| Texnologiya | Vazifasi |
| --- | --- |
| **Python 3.12 / 3.11** | Asosiy til (Docker 3.11, README 3.12) |
| **FastAPI (Async)** | Web framework, async endpoints |
| **Uvicorn** | ASGI server |
| **SQLAlchemy (async)** | ORM, `asyncpg` driver bilan |
| **Alembic** | DB migration |
| **PostgreSQL** | Asosiy ma'lumotlar bazasi |
| **Redis** | Cache, token blacklist, distributed lock, scheduler jobstore |
| **Kafka (aiokafka)** | Event publishing (xabar yuborish) |
| **MinIO** | File storage (S3-ga o'xshash) |
| **APScheduler** | Rejalashtirilgan vazifalar (cron-jobs) |
| **SQLAdmin / starlette_admin** | Admin panel (dashboard) |
| **Pydantic v2 / pydantic-settings** | Validatsiya va konfiguratsiya |
| **httpx** | Tashqi servislarga HTTP so'rovlar |
| **Docker** | Konteynerizatsiya |

---

# 1. Loyiha Strukturasi (Mermaid)

```mermaid
graph TD
    Root["fastapi-structure/"]

    Root --> Docker["Dockerfile / entrypoint.sh"]
    Root --> Make["Makefile"]
    Root --> Alembic["alembic.ini"]
    Root --> App["app/"]

    App --> Main["main.py — kirish nuqtasi"]
    App --> Core["core/ — infratuzilma"]
    App --> Api["api/v1/ — routerlar (endpointlar)"]
    App --> Deps["dependencies/ — DI factory'lar"]
    App --> Services["services/ — biznes logika"]
    App --> Repos["repositories/ — DB so'rovlari"]
    App --> Models["models/ — SQLAlchemy jadvallar"]
    App --> Schemas["schemas/ — Pydantic DTO"]
    App --> Uow["uow/ — Unit of Work"]
    App --> Clients["clients/ — tashqi servis adapterlari"]
    App --> Admin["admin/ — admin panel"]
    App --> Db["db/ — engine va session"]
    App --> Migrations["migrations/ — Alembic"]
    App --> Utils["utils/ — bot, http_client"]
```

## Qatlamli arxitektura oqimi (eng muhim diagramma)

Bu loyiha **Clean Architecture / qatlamli (layered)** uslubda yozilgan. So'rov quyidagi yo'lni bosib o'tadi:

```mermaid
flowchart LR
    Client([Client / Frontend]) -->|HTTP| Router["api/v1/*.py<br/>(Router)"]
    Router -->|Depends| Dep["dependencies/*.py<br/>(Factory)"]
    Dep -->|yaratadi| Service["services/*.py<br/>(Biznes logika)"]
    Service -->|ishlatadi| Uow["uow/impl.py<br/>(Unit of Work)"]
    Service -->|ishlatadi| Cli["clients/*.py<br/>(SSO, Gateway)"]
    Uow -->|saqlaydi| Repo["repositories/*.py<br/>(Repository)"]
    Repo -->|SQL| Model["models/*.py<br/>(SQLAlchemy)"]
    Model -->|asyncpg| DB[(PostgreSQL)]
    Cli -->|httpx| Ext([Tashqi servislar])

    Service -.javob.-> Schema["schemas/*.py<br/>(Pydantic)"]
    Schema -.JSON.-> Client
```

**Asosiy qoida:** har bir qatlam faqat **pastdagi** qatlamni biladi. Router servisni biladi, servis repository'ni biladi, lekin repository router haqida hech narsa bilmaydi. Bu loyihani test qilish va o'zgartirishni osonlashtiradi.

---

# 2. Ildizdagi (root) fayllar

## `Dockerfile`
Konteyner qurish retsepti:

```dockerfile
FROM harbor.mf.uz/global/library/python:3.11-slim-bookworm
WORKDIR /app
ENV TZ=Asia/Tashkent           # Toshkent vaqt zonasi
RUN apt-get install ... netcat-openbsd libpq-dev gcc ...
COPY ./requirements.txt /app/requirements.txt
RUN pip install -r requirements.txt
COPY . /app
RUN python app/collect_static.py   # admin panel static fayllarini yig'adi
EXPOSE 9080
HEALTHCHECK ... CMD curl -f http://localhost:9080/health   # sog'liqni tekshirish
ENTRYPOINT ["/usr/src/app/entrypoint.sh"]
CMD ["python", "app/main.py"]
```

Muhim nuqtalar:
- `netcat-openbsd` — `entrypoint.sh` da postgres ishga tushishini kutish uchun ishlatiladi
- `libpq-dev`, `gcc` — `asyncpg`/`psycopg` kompilyatsiyasi uchun
- `HEALTHCHECK` — `/health` endpointni har 10 soniyada tekshiradi

## `entrypoint.sh`
Konteyner ishga tushganda **birinchi** bajariladigan skript:

```sh
if [ "$DATABASE" = "postgres" ]; then
    while ! nc -z $SQL_HOST $SQL_PORT; do sleep 0.1; done   # postgres tayyor bo'lguncha kut
fi
alembic upgrade head    # migratsiyalarni avtomatik qo'llaydi
exec "$@"               # keyin CMD ni (python app/main.py) ishga tushiradi
```

Bu pattern "wait-for-it" deb ataladi — DB tayyor bo'lmaguncha ilova ishga tushmaydi.

## `Makefile`
Tez-tez ishlatiladigan buyruqlar qisqartmasi:

```makefile
install:          # pip3 install -r requirements.txt
lint:             # isort + black (kodni formatlash)
make_migrations:  # alembic revision --autogenerate (yangi migratsiya)
migrate:          # alembic upgrade head (migratsiyani qo'llash)
run:              # python3 app/main.py
deps-up/down:     # docker compose up/down
```

## `alembic.ini`
Alembic (DB migration vositasi) konfiguratsiyasi:
- `script_location = app/migrations` — migratsiyalar qayerda turishi
- `prepend_sys_path = . app` — `app` papkasini import yo'liga qo'shadi (shuning uchun kodda `from core.config` deb yozish mumkin, `from app.core.config` emas)

---

# 3. `app/core/` — Infratuzilma qatlami

Bu papka ilovaning "asboblar qutisi": konfiguratsiya, xavfsizlik, tashqi tizimlar ulanishlari.

## `core/config.py` — Markazlashtirilgan konfiguratsiya

Pydantic `BaseSettings` yordamida barcha sozlamalar `.env` fayldan o'qiladi va **tipi tekshiriladi**:

```python
class Settings(BaseSettings):
    POSTGRES_PORT: int        # agar .env da raqam bo'lmasa, ilova ishga tushmaydi
    SECRET_KEY: str
    REDIS_URL: str
    ...

    @property
    def DATABASE_URL(self):   # computed property — boshqa qiymatlardan quriladi
        return f"postgresql+asyncpg://{self.POSTGRES_USER}:{self.POSTGRES_PASSWORD}@..."

    model_config = SettingsConfigDict(env_file=".env", extra="ignore")

settings = Settings()  # bir marta yaratiladi, hamma joyda import qilinadi (singleton)
```

**Nima uchun yaxshi:** sozlamalar bir joyda, tipi tekshirilgan. Agar `.env` da `POSTGRES_PORT` yo'q bo'lsa, ilova xato bilan to'xtaydi — production'da noto'g'ri konfiguratsiya bilan ishlamaydi.

## `core/exceptions.py` — Maxsus xato turlari

Barcha biznes xatolar uchun ierarxiya:

```mermaid
classDiagram
    Exception <|-- AppException
    AppException <|-- BusinessException
    AppException <|-- NotFoundException
    AppException <|-- ExternalServiceException
    AppException <|-- CustomValidationException

    class AppException {
        +str error
        +str detail
        +int status_code
        +str stack_trace
    }
    class BusinessException {
        status_code = 400
    }
    class NotFoundException {
        status_code = 404
    }
    class ExternalServiceException {
        status_code = 502
    }
    class CustomValidationException {
        status_code = 422
    }
```

```python
class NotFoundException(AppException):
    def __init__(self, detail: str):
        super().__init__(error="NotFound", detail=detail, status_code=404)
```

**Foydasi:** kodda `raise NotFoundException("Xodim topilmadi")` deb yozsangiz, bu avtomatik 404 status va to'g'ri JSON formatga aylanadi. HTTP status kodlarini har joyda qo'lda yozish shart emas.

`AppException.__init__` ichida `traceback.format_stack()` chaqiriladi — xato qayerda ko'tarilganini eslab qoladi (debug uchun).

## `core/error_handlers.py` — Global xato ushlovchilar

Bu yerda 4 ta handler bor, ular `main.py` da ro'yxatdan o'tkaziladi:

```mermaid
flowchart TD
    Raise["Kodda exception ko'tariladi"] --> Type{Qanday tur?}
    Type -->|AppException| H1["app_exception_handler<br/>→ ErrorResponse JSON"]
    Type -->|RequestValidationError| H2["validation_exception_handler<br/>→ maydon xatolari ro'yxati"]
    Type -->|HTTPException| H3["http_exception_handler<br/>→ 5xx bo'lsa Telegram'ga yuboradi"]
    Type -->|boshqa har qanday| H4["unhandled_exception_handler<br/>→ 500 + Telegram'ga to'liq traceback"]
```

Eng muhim qism — `unhandled_exception_handler`: production'da kutilmagan xato bo'lsa, **Telegram botga** to'liq traceback yuboradi (3500 belgidan oshsa ikki xabarga bo'lib). Foydalanuvchiga esa faqat umumiy "Kutilmagan xatolik" xabari qaytadi.

```python
async def app_exception_handler(request: Request, exc: AppException):
    return JSONResponse(
        status_code=exc.status_code,
        content=ErrorResponse(status_code=exc.status_code, error=exc.error, detail=exc.detail).model_dump(),
    )
```

## `core/security.py` — JWT autentifikatsiya

Token yaratish, tekshirish va logout:

```mermaid
sequenceDiagram
    participant U as User
    participant S as create_access_token
    participant V as verify_access_token
    participant R as Redis
    participant L as logout

    U->>S: userId, userPinfl
    S-->>U: JWT (jti, exp bilan)
    U->>V: JWT yuboradi
    V->>R: blacklist:{jti} bormi?
    R-->>V: yo'q
    V-->>U: payload (ruxsat)
    U->>L: logout qiladi
    L->>R: setex blacklist:{jti} (TTL=exp)
    Note over R: endi bu token bloklangan
```

```python
def create_access_token(userId, userPinfl):
    jti = str(uuid.uuid4())       # har token uchun noyob ID
    payload = {"user_id": ..., "jti": jti, "exp": expire}
    return jwt.encode(payload, SECRET_KEY, algorithm="HS256")
```

**Aqlli yechim — logout:** JWT odatda server tomonida bekor qilinmaydi (stateless). Bu yerda har tokenga `jti` (JWT ID) beriladi. Logout qilganda `jti` Redis'ga "blacklist" sifatida saqlanadi, TTL esa token amal qilish vaqtiga teng. `verify_access_token` har safar blacklist'ni tekshiradi.

## `core/redis.py` — Redis ulanishi

```python
redis_client = redis.Redis.from_url(settings.REDIS_URL)   # async Redis
# APScheduler joblar Redis'da saqlanadi (server restartdan keyin yo'qolmaydi)
jobstores = {"default": RedisJobStore(...)}
scheduler = AsyncIOScheduler(jobstores=jobstores)
```

## `core/redis_lock.py` — Distributed Lock

Ko'p instansiyali (bir nechta konteyner) deploymentda **rejalashtirilgan vazifa faqat bitta instansiyada** bajarilishi kerak. Bu Redis lock buni ta'minlaydi:

```mermaid
flowchart TD
    A["Instansiya 1: SET lock NX EX"] -->|muvaffaqiyat| B["Vazifani bajaradi"]
    C["Instansiya 2: SET lock NX EX"] -->|key allaqachon bor| D["Vazifani o'tkazib yuboradi"]
    B --> E["Lua script bilan lock'ni bo'shatadi"]
```

```python
result = await self.redis_client.set(
    self.lock_name, self.lock_value,
    nx=True,           # faqat agar mavjud bo'lmasa o'rnatadi (atomar)
    ex=self.timeout,   # timeout dan keyin o'zi o'chadi (deadlock bo'lmaydi)
)
```

`release()` da **Lua script** ishlatiladi — faqat lock egasi o'chira oladi (boshqa instansiyaning lock'ini xato o'chirib yubormaslik uchun):

```lua
if redis.call("get", KEYS[1]) == ARGV[1] then
    return redis.call("del", KEYS[1])
else return 0 end
```

`__aenter__`/`__aexit__` mavjud, ya'ni `async with RedisDistributedLock(...)` sifatida ishlatish mumkin.

## `core/kafka_publisher.py` — Event publishing

Abstract klass + konkret implementatsiya (Strategy/Adapter pattern):

```mermaid
classDiagram
    MessagePublisher <|-- KafkaPublisher
    class MessagePublisher {
        <<abstract>>
        +publish(topic, type, created_at, payload, aggregate_id)
        +check_kafka_cluster_health()
    }
    class KafkaPublisher {
        -AIOKafkaProducer _producer
        +start()
        +stop()
        +publish(...)
    }
```

```python
async def publish(self, topic, type, created_at, payload, aggregate_id):
    headers = [("eventType", type.encode()), ("recordTime", ...), ("id", uuid...)]
    await self._producer.send_and_wait(
        topic=topic, headers=headers,
        value=json.dumps(payload).encode(),
        key=aggregate_id.encode(),   # bir aggregate xabarlar tartibini saqlaydi
    )
```

**Nima uchun abstract:** `MessagePublisher` interfeysi tufayli ertaga Kafka o'rniga RabbitMQ ishlatmoqchi bo'lsangiz, faqat yangi klass yozasiz, qolgan kod o'zgarmaydi (Dependency Inversion).

## `core/minio.py` — File storage

S3-ga mos MinIO bilan ishlash. `FileStorage` abstract klass + `MinioClient`:

```python
class FileStorage(ABC):
    @abstractmethod
    async def upload_file_from_buffer(...): pass
    @abstractmethod
    async def download_file(...): pass
```

Muhim: MinIO kutubxonasi **sinxron**, lekin ilova **async**. Shuning uchun bloklovchi chaqiruvlar `anyio.to_thread.run_sync(...)` orqali alohida thread'da bajariladi — event loop bloklanmaydi.

## `core/permissions.py` — Rol asosida ruxsat

```python
def get_hr_user(user: RequestUserDep) -> User:
    if RoleKey.HR == user.currentRole.key:
        return user
    raise HTTPException(403, "Request user is not HR")

HRUserDep = Annotated[User, Depends(get_hr_user)]
```

`HRUserDep` ni endpoint'ga qo'shsangiz, faqat HR roli bilan kira oladi. (Eslatma: bu fayl `models.role` ga bog'liq, lekin `models/role.py` repozitoriyada yo'q — bu template'ning kengaytirilishi kutiladigan qismi.)

---

# 4. `app/db/` — Ma'lumotlar bazasi ulanishi

## `db/session.py`

```python
engine = create_async_engine(url=settings.DATABASE_URL, echo=settings.DEBUG)
session_factory = async_sessionmaker(engine)   # session "fabrikasi"

async def get_db() -> AsyncGenerator[AsyncSession, None]:
    async with session_factory() as session:
        try:
            yield session       # FastAPI dependency sifatida ishlatiladi
        finally:
            await session.close()
```

- `engine` — DB bilan ulanishlar puli (bir marta yaratiladi)
- `session_factory` — har so'rov uchun yangi session yaratadi
- `echo=settings.DEBUG` — DEBUG rejimida barcha SQL so'rovlar log'ga chiqadi

---

# 5. `app/models/` — SQLAlchemy modellar (jadvallar)

Modellar DB jadvallarini Python klasslari sifatida ifodalaydi.

## `models/base.py` — Barcha modellar uchun asos

```python
class Base(DeclarativeBase):
    id: Mapped[int] = mapped_column(primary_key=True, autoincrement=True)
    createdAt: Mapped[datetime] = mapped_column(..., name="created_at")
    updatedAt: Mapped[datetime] = mapped_column(..., onupdate=..., name="updated_at")
    deletedAt: Mapped[datetime | None] = mapped_column(..., name="deleted_at", index=True)
```

**Muhim patternlar:**
- Har bir jadval avtomatik `id`, `created_at`, `updated_at`, `deleted_at` oladi
- `deletedAt` — **soft delete** uchun (yozuvni o'chirmasdan, sana qo'yiladi)
- Python'da `createdAt` (camelCase), DB'da `created_at` (snake_case) — `name=` orqali bog'lanadi

## `models/user.py` va `models/organization.py`

```python
class User(Base):
    __tablename__ = "users"
    pinfl: Mapped[str] = mapped_column("pinfl", String, nullable=False)
    userUuid: Mapped[uuid.UUID] = mapped_column(default=uuid.uuid4, unique=True, index=True)
    gender: Mapped[GenderType] = mapped_column(Enum(GenderType), nullable=True)
    ...
    def __str__(self):
        return self.fullName   # admin panelда ko'rsatish uchun
```

```mermaid
erDiagram
    USERS {
        int id PK
        string username UK
        uuid user_uuid UK
        string pinfl
        string full_name
        enum gender
        timestamp created_at
        timestamp deleted_at
    }
    ORGANIZATIONS {
        int id PK
        string name_uz
        string tin UK
        string director_name
        enum status
        timestamp created_at
        timestamp deleted_at
    }
```

`enum.Enum` ishlatilgan: `GenderType` (Male/Female), `OrgStatus` (ACTIVE/NOT_ACTIVE) — DB darajasida cheklash qo'yiladi.

---

# 6. `app/schemas/` — Pydantic DTO (ma'lumot uzatish obyektlari)

Model — DB jadvali. Schema — **API kirish/chiqish formati**. Ularni ajratish muhim, chunki client'ga modelning barcha maydonlarini (masalan, `password`) ko'rsatmaslik kerak.

```mermaid
flowchart LR
    Req["Request JSON"] -->|validatsiya| InSchema["UserOneIdRequest"]
    InSchema --> Service["Service"]
    Service --> Model["User (DB model)"]
    Model -->|from_attributes| OutSchema["OrganizationResponse"]
    OutSchema -->|model_dump| Resp["Response JSON"]
```

## `schemas/organization.py`

```python
class OrganizationResponse(BaseModel):
    id: int
    name_uz: str
    tin: str
    status: OrgStatus
    model_config = ConfigDict(from_attributes=True)  # SQLAlchemy obyektidan to'g'ridan-to'g'ri o'qiydi
```

`from_attributes=True` muhim — bu SQLAlchemy `Organization` obyektini avtomatik `OrganizationResponse` ga aylantirishga imkon beradi.

## `schemas/error.py` — standartlashtirilgan xato formati

```python
class ErrorResponse(BaseModel):
    status_code: int
    error: str
    detail: Optional[str] = None

class ValidationErrorResponse(BaseModel):
    status_code: int
    error: str
    detail: list[FieldError]   # qaysi maydonda qanday xato
```

Barcha xatolar bir xil JSON ko'rinishga ega bo'ladi — frontend ishi osonlashadi.

## `schemas/sso.py` va `schemas/gateway.py`

Tashqi servislardan keladigan ma'lumot formatlari. Masalan `SsoUserInfoResponse` (One ID dan keluvchi xodim ma'lumotlari), `UserJobDetailResponse`, `UserPassportInfoResponse` (Gateway'dan).

---

# 7. `app/repositories/` — Repository pattern

Repository — **faqat** DB bilan ishlaydigan qatlam. Biznes logika yo'q, faqat CRUD.

```python
class UserRepository:
    def __init__(self, session: AsyncSession, model=User):
        self.session = session
        self.model = model

    async def add(self, user: User) -> User:
        self.session.add(user)
        await self.session.flush()       # INSERT ni bajaradi, lekin commit qilmaydi
        await self.session.refresh(user) # DB dan id va default qiymatlarni oladi
        return user

    async def get_user_by_pinfl(self, pinfl: str) -> User | None:
        stmt = select(User).where(User.pinfl == pinfl)
        result = await self.session.execute(stmt)
        return result.scalar_one_or_none()
```

**Muhim:** `flush()` qiladi, lekin `commit()` qilmaydi. Commit qarorini **Unit of Work** o'ziga oladi. Bu transaction'ni boshqarishni bir joyga to'playdi.

---

# 8. `app/uow/` — Unit of Work pattern

Bu loyihaning eng muhim patternlaridan biri. Unit of Work bitta biznes operatsiyasini **bitta transaction** ichida bajaradi.

```mermaid
sequenceDiagram
    participant Svc as Service
    participant Uow as UnitOfWork
    participant Sess as AsyncSession
    participant DB as PostgreSQL

    Svc->>Uow: async with uow
    Uow->>Sess: __aenter__: session yaratadi
    Uow->>Uow: repositorylarni init qiladi (user, organization)
    Svc->>Uow: uow.user.add(...)
    Uow->>DB: flush (INSERT, commit emas)
    alt xato bo'lmasa
        Uow->>DB: __aexit__: commit
    else xato bo'lsa
        Uow->>DB: __aexit__: rollback
    end
    Uow->>Sess: session.close()
```

## `uow/ports.py` — Interfeys (abstraksiya)

```python
class UnitOfWorkABC(ABC):
    @abstractmethod
    async def __aenter__(self): ...
    @abstractmethod
    async def __aexit__(self, exc_type, exc_val, exc_tb): ...
    @abstractmethod
    async def commit(self): ...
    @abstractmethod
    async def rollback(self): ...
```

## `uow/impl.py` — Implementatsiya

```python
class UnitOfWork(UnitOfWorkABC):
    async def __aenter__(self) -> Self:
        self._session = self._session_factory()
        self.user = UserRepository(self._session, User)            # repositorylar
        self.organization = OrganizationRepository(self._session, Organization)
        return self

    async def __aexit__(self, exc_type, exc_value, traceback):
        if exc_type:
            await self.rollback()   # xato bo'lsa hammasini bekor qiladi
        else:
            await self.commit()     # hammasi yaxshi bo'lsa saqlaydi
```

**Foydasi:** Service ichida `async with self.uow:` deb yozsangiz, ichidagi barcha DB amallari yo birga saqlanadi (commit), yoki birga bekor qilinadi (rollback). Bir nechta jadvalni o'zgartirishda ma'lumot izchilligi (consistency) ta'minlanadi.

---

# 9. `app/clients/` — Tashqi servis adapterlari

Tashqi tizimlar (One ID SSO, Gateway) bilan ishlash. Bu yerda ham **abstract klass (port) + implementatsiya** ishlatiladi.

```mermaid
classDiagram
    SsoClientABC <|-- SsoClient
    GatewayClientABC <|-- GatewayClient
    class SsoClientABC {
        <<abstract>>
        +get_access_token(code, codeVerify, redirectUri)
        +get_user_info(accessToken)
    }
    class GatewayClientABC {
        <<abstract>>
        +get_access_token(...)
        +get_user_job_detail(...)
        +get_user_passport_info(...)
        +get_org_account_numbers(...)
        +fetch_org_district(...)
    }
```

## `clients/sso.py` — One ID integratsiyasi

```python
class SsoClient(SsoClientABC):
    async def get_access_token(self, code, codeVerify, redirectUri) -> SsoSuccessResponse:
        self.url = f"{self.base_url}/oauth2/token?...&code={code}..."
        try:
            response = await self.client.post(url=self.url, ..., auth=self.auth)
        except httpx.RequestError as e:
            raise ExternalServiceException(service="SSO", detail=f"...: {e}")
        if response.status_code != 200:
            raise CustomValidationException(detail=f"...: {response.text}")
        return SsoSuccessResponse(accessToken=response.json().get("access_token"))
```

**Patternlar:**
- Ulanish xatosi → `ExternalServiceException` (502)
- Yomon javob → `CustomValidationException` (422)
- Tashqi servis JSON'i ichki schema'ga (Pydantic) o'tkaziladi — qolgan kod tashqi format detallarini bilmaydi

## `clients/gateway.py` — Gateway integratsiyasi

Eng katta client. Xodim ish joyi, pasport ma'lumoti, tashkilot hisob raqamlari, tuman kodi kabilarni oladi. Tashqi JSON'ni o'zining `schemas/gateway.py` modellariga aylantiradi (anti-corruption layer):

```python
userPositions = [
    UserJobPositionDetail(
        orgTin=position["org_tin"],
        positionName=position["position"],
        ...
    ) for position in positions
]
```

---

# 10. `app/services/` — Biznes logika qatlami

Service — **biznes qoidalari** joylashgan joy. U Repository, Client va UoW'ni birlashtiradi.

## `services/user.py`

```python
class UserService:
    def __init__(self, uow: UnitOfWorkABC, ssoClient: SsoClient):
        self.uow = uow
        self.ssoClient = ssoClient

    async def one_id(self, code, codeVerify, redirectUri) -> UserOneIdResponse:
        ssoRes = await self.ssoClient.get_access_token(code, codeVerify, redirectUri)
        userInfo = await self.ssoClient.get_user_info(ssoRes.accessToken)

        async with self.uow:
            user = await self.uow.user.get_user_by_pinfl(userInfo.pinfl)
            if not user:
                raise NotFoundException(detail="Xodim tizimda mavjud emas")
            token = create_access_token(userId=user.id, userPinfl=user.pinfl)

        return UserOneIdResponse(accessToken=token, ssoToken=ssoRes.accessToken)
```

`one_id` oqimi (login):

```mermaid
sequenceDiagram
    participant FE as Frontend
    participant Svc as UserService
    participant SSO as SsoClient (One ID)
    participant Uow as UnitOfWork
    participant DB as PostgreSQL

    FE->>Svc: one_id(code, codeVerify, redirectUri)
    Svc->>SSO: get_access_token(...)
    SSO-->>Svc: ssoToken
    Svc->>SSO: get_user_info(ssoToken)
    SSO-->>Svc: userInfo (pinfl, ism...)
    Svc->>Uow: get_user_by_pinfl(pinfl)
    Uow->>DB: SELECT
    alt topilmadi
        Svc-->>FE: 404 NotFoundException
    else topildi
        Svc->>Svc: create_access_token(user)
        Svc-->>FE: accessToken + ssoToken
    end
```

`services/organization.py` ham xuddi shu strukturada — `get_org_by_tin` TIN bo'yicha tashkilot qidiradi, topilmasa `NotFoundException`.

> Eslatma: `dependencies/` da `services/out.py`, `services/region.py`, `services/role.py`, `services/org_type.py`, `services/district.py` ga havola bor, lekin bu fayllar template'da hali yozilmagan — kelajakda kengaytirish uchun "joy qoldirilgan".

---

# 11. `app/dependencies/` — Dependency Injection (DI)

FastAPI'ning eng kuchli xususiyati — **Dependency Injection**. Bu papka "factory" funksiyalarini saqlaydi: har bir service/client'ni qanday yaratishni bilgan funksiyalar.

```mermaid
flowchart TD
    Endpoint["@router.post: service: UserServiceDep"]
    Endpoint --> F1["get_user_service"]
    F1 --> F2["get_unit_of_work (UowDependency)"]
    F1 --> F3["get_sso_client (SsoClientDep)"]
    F2 --> Uow["UnitOfWork(session_factory)"]
    F3 --> Sso["SsoClient()"]
    F1 --> Result["UserService(uow, ssoClient)"]
```

## `dependencies/uow.py` — UoW provayderi

```python
async def get_unit_of_work() -> AsyncGenerator[UnitOfWorkABC, None]:
    async with UnitOfWork(session_factory) as uow:
        yield uow

UowDependency = Annotated[UnitOfWorkABC, Depends(get_unit_of_work)]
```

Fayldagi izohda muhim dars bor: avval `return UnitOfWork(...)` qilingan (xato) — `__aenter__` chaqirilmagan, session ochilmagan, repositorylar init bo'lmagan. To'g'ri yechim — `async with` orqali kontekstga kirish.

## `dependencies/user.py` — Eng to'liq misol

Bu yerda bir nechta DI mavjud:

```python
# 1. Service factory
def get_user_service(uow: UowDependency, ssoClient: SsoClientDep):
    return UserService(uow=uow, ssoClient=ssoClient)
UserServiceDep = Annotated[UserService, Depends(get_user_service)]

# 2. Joriy foydalanuvchini token'dan olish
async def get_request_user(uow, credentials=Depends(HTTPBearer())):
    token = credentials.credentials
    payload = await verify_access_token(token)
    user = await uow.user.get_user_by_id(int(payload["user_id"]))
    return user
RequestUserDep = Annotated[User, Depends(get_request_user)]

# 3. Internal servislar uchun Basic Auth
async def get_basic_auth_user(credentials=Depends(HTTPBasic())):
    if not (secrets.compare_digest(...) and secrets.compare_digest(...)):
        raise HTTPException(401, ...)
    return credentials.username
BasicAuthUserDep = Annotated[str, Depends(get_basic_auth_user)]
```

`secrets.compare_digest` — parolni vaqt-hujumiga (timing attack) qarshi taqqoslaydi.

## DI patternining go'zalligi

`Annotated[Type, Depends(factory)]` ni alias sifatida saqlash — bu loyihaning asosiy uslubi. Endpoint'da shunchaki:

```python
@router.get("/{pinfl}")
async def get_user(pinfl: str, service: UserServiceDep):   # FastAPI o'zi service yaratadi
    return await service.get_user_by_pinfl(pinfl)
```

FastAPI avtomatik ravishda: session ochadi → repositorylar yaratadi → UoW yasaydi → SSO client yasaydi → UserService yasaydi → endpoint'ga beradi. So'rov tugagach hammasini tozalaydi.

Boshqa dependency'lar (`org.py`, `sso.py`, `gateway.py`, `region.py`, `district.py`, `role.py`, `org_type.py`, `outbox.py`) ham aynan shu naqsh bo'yicha yozilgan.

---

# 12. `app/api/v1/` — Routerlar (Endpointlar)

Bu yerda HTTP endpointlar e'lon qilinadi. Diqqat qiling — endpointlar **juda yupqa** (thin): ular faqat service'ni chaqiradi, biznes logika ularda yo'q.

## `api/v1/user.py`

```python
router = APIRouter(prefix="/user", tags=["User"])

@router.get("/{pinfl}")
async def get_user_by_pinfl(pinfl: str, service: UserServiceDep):
    return await service.get_user_by_pinfl(pinfl)

@router.post("/oneId/", response_model=UserOneIdResponse)
async def one_id(params: UserOneIdRequest, service: UserServiceDep) -> UserOneIdResponse:
    return await service.one_id(code=params.code, codeVerify=params.codeVerify, redirectUri=params.redirectUri)
```

- `APIRouter(prefix=..., tags=...)` — endpointlarni guruhlaydi (Swagger'da tab bo'ladi)
- `response_model` — javob avtomatik shu schema'ga moslab filtrlanadi
- `params: UserOneIdRequest` — kirish JSON avtomatik validatsiya qilinadi

## `api/v1/organization.py`

```python
router = APIRouter(prefix="/organizations", tags=["Organization"])

@router.get("/{tin}", response_model=OrganizationResponse)
async def get_organization_by_tin(tin: str, service: OrgServiceDep):
    return await service.get_org_by_tin(tin)
```

Routerlar `main.py` da ro'yxatdan o'tkaziladi:

```python
all_routers = [user_router, organization_router]
for router in all_routers:
    app.include_router(router, prefix="/api/v1/uzassets")
# Yakuniy URL: /api/v1/uzassets/user/{pinfl}
```

---

# 13. `app/admin/` — Admin panel (Dashboard)

`starlette_admin` yordamida avtomatik admin panel quriladi.

```mermaid
flowchart TD
    Setup["setup_admin(app)"] --> Auth["AdminAuth — login/logout"]
    Setup --> V1["UserAdminView(User)"]
    Setup --> V2["OrganizationAdminView(Organization)"]
    V1 --> Base["BaseAdminView"]
    V2 --> Base
    Base --> Mount["admin.mount_to(app)"]
    Mount --> URL["/api/v1/uzassets/dashboard/"]
```

## `admin/auth.py` — Admin autentifikatsiyasi

```python
class AdminAuth(AuthProvider):
    async def login(self, username, password, ...):
        if username == settings.ADMIN_USERNAME and password == settings.ADMIN_PASSWORD:
            request.session.update({"username": username})  # session'ga saqlaydi
            return response
        raise LoginFailed("Invalid username or password")
```

## `admin/views/base.py` — Umumiy view sozlamalari

```python
class BaseAdminView(ModelView):
    exclude_fields_from_create = ["createdAt", "updatedAt"]  # bu maydonlar ko'rsatilmaydi
    def can_delete(self, request) -> bool:
        if settings.APP_MODE == "PRODUCTION":
            return False   # production'da o'chirish taqiqlangan
        return True
```

`admin/views/user.py` va `organization.py` — har bir model uchun konkret view (nom, saralash maydonlari).

## `admin/__init__.py` — HTTPS muammosini hal qilish

Nginx ortida ishlaganda admin panel `http://` URL'lar generatsiya qilib qo'yardi. `HTTPSAdmin` va `HTTPSStaticFiles` klasslari URL'larni majburan `https://` ga aylantiradi.

---

# 14. `app/migrations/` — Alembic migratsiyalari

DB sxemasini versiyalash (har o'zgarish alohida fayl).

## `migrations/env.py`

```python
target_metadata = Base.metadata   # barcha modellar metadata'si
config.set_main_option("sqlalchemy.url", settings.DATABASE_URL)
# async engine bilan migratsiyalarni bajaradi
```

Muhim: `from models.user import User` va `from models.organization import Organization` — modellarni import qilish **shart**, aks holda Alembic ularni "ko'rmaydi" va `--autogenerate` jadval yaratmaydi.

## `migrations/versions/`

```mermaid
flowchart LR
    Base0["(boshlang'ich)"] --> M1["7df8fe487462<br/>users jadvali"]
    M1 --> M2["fab628a152d2<br/>organizations jadvali"]
```

Har bir fayl `upgrade()` (jadval yaratish) va `downgrade()` (orqaga qaytarish) funksiyalariga ega. `down_revision` orqali zanjir tuziladi.

---

# 15. `app/main.py` — Ilovaning yuragi

Bu fayl hamma narsani birlashtiradi.

## `lifespan` — Startup va Shutdown

```mermaid
stateDiagram-v2
    [*] --> Startup
    Startup --> SetupLogging: setup_logging()
    SetupLogging --> SetupAdmin: setup_admin(app)
    SetupAdmin --> RedisPing: redis.ping()
    RedisPing --> KafkaStart: kafka_publisher.start()
    KafkaStart --> Scheduler: scheduler.start()
    Scheduler --> Running: yield (ilova ishlaydi)
    Running --> Shutdown: SIGTERM
    Shutdown --> StopScheduler: scheduler.shutdown()
    StopScheduler --> StopKafka: kafka.stop()
    StopKafka --> [*]
```

```python
@asynccontextmanager
async def lifespan(app: FastAPI):
    setup_logging()
    setup_admin(app)
    await redis_client.ping()          # Redis majburiy
    await kafka_publisher.start()      # Kafka producer ishga tushadi
    app.state.publisher = kafka_publisher  # boshqa joylarda ishlatish uchun
    scheduler.start()
    yield                              # <-- ilova shu yerda ishlaydi
    # quyidagilar shutdown'da bajariladi:
    await scheduler.shutdown(wait=False)
    await kafka_publisher.stop()
```

`yield` dan **yuqorisi** — startup, **pasti** — shutdown. Resurslarni to'g'ri ochish/yopish ta'minlanadi.

## Middleware'lar

```mermaid
flowchart TD
    Req["So'rov keladi"] --> M1["TrustedHostMiddleware<br/>(faqat ruxsat etilgan hostlar)"]
    M1 --> M2["CORSMiddleware<br/>(cross-origin so'rovlar)"]
    M2 --> M3["SessionMiddleware<br/>(admin sessiyasi)"]
    M3 --> M4["https_middleware<br/>(HTTPS + xavfsizlik header'lari)"]
    M4 --> Route["Endpoint"]
    Route --> Resp["Javob"]
```

```python
@app.middleware("http")
async def https_middleware(request, call_next):
    # x-forwarded-proto headerni tekshiradi (Nginx ortida HTTPS aniqlash)
    response = await call_next(request)
    response.headers["Strict-Transport-Security"] = "max-age=31536000; ..."
    response.headers["Content-Security-Policy"] = "..."   # XSS himoyasi
    return response
```

- **TrustedHostMiddleware** — Host header hujumlaridan himoya
- **CORSMiddleware** — frontend boshqa domendan so'rov yuborishi uchun
- **SessionMiddleware** — admin panel login sessiyasi uchun
- **https_middleware** — proxy ortidagi HTTPS va xavfsizlik header'lari

## Exception handler'lar va Health check

```python
app.add_exception_handler(AppException, app_exception_handler)
app.add_exception_handler(RequestValidationError, validation_exception_handler)
app.add_exception_handler(StarletteHTTPException, http_exception_handler)
app.add_exception_handler(Exception, unhandled_exception_handler)

@app.get("/health")   # Docker HEALTHCHECK uchun
async def health_check():
    # Kafka, Redis, PostgreSQL tekshiriladi; biror narsa ishlamasa 503 qaytaradi
```

## OpenAPI sozlamalari

```python
app = FastAPI(
    title="UzAssets",
    docs_url="/api/v1/uzassets/swagger/",       # Swagger UI manzili
    openapi_url="/api/v1/uzassets/openapi.json",
    redoc_url=None,                              # ReDoc o'chirilgan
    lifespan=lifespan,
)
```

---

# 16. `app/logging_config.py` va `app/collect_static.py`

## `logging_config.py`
Uvicorn-uslubidagi ranglangan loglar. Shovqinli kutubxonalar (`aiokafka`, `sqlalchemy.engine`) ERROR darajasiga tushiriladi:

```python
"loggers": {
    "aiokafka": {"level": "ERROR"},          # ortiqcha INFO loglarni o'chiradi
    "sqlalchemy.engine": {"level": "ERROR"},
}
```

## `collect_static.py`
Docker build paytida admin panel (sqladmin, starlette_admin) static fayllarini `/app/statics` ga ko'chiradi. Dockerfile'da `RUN python app/collect_static.py` chaqiriladi.

---

# 17. `app/utils/` — Yordamchi vositalar

`utils/__init__.py` faqat `from .telegram_bot import bot` qatorini o'z ichiga oladi. `telegram_bot.py` va `http_client.py` fayllari `.gitignore` orqali repozitoriyaga kiritilmagan (sirli token bo'lgani uchun), lekin kodda ishlatiladi:
- `utils.bot` — xatolarni Telegram'ga yuborish (`error_handlers.py` da)
- `utils.http_client.get_http_client()` — umumiy `httpx` client (`clients/` da)

---

# 18. Asosiy konsepsiyalar — xulosa

## Routing
- `APIRouter` bilan endpointlar guruhlanadi (`prefix`, `tags`)
- `main.py` da `include_router` orqali ulanadi
- Path parametr (`/{pinfl}`), Body (`UserOneIdRequest`), `response_model`

## Middleware
- So'rov/javobni qayta ishlash zanjiri
- Xavfsizlik (HTTPS, CSP, TrustedHost), CORS, Session

## Dependency Injection
- `Annotated[Type, Depends(factory)]` — kod takrorlanmaydi
- FastAPI obyektlar daraxtini o'zi quradi va tozalaydi
- Test qilishda mock berish oson

## Layered Architecture
```mermaid
mindmap
  root((FastAPI Template))
    Routing
      APIRouter
      response_model
      tags / prefix
    Qatlamlar
      Router
      Service
      Repository
      Model
    Patternlar
      Unit of Work
      Repository
      Dependency Injection
      Abstract Port + Impl
    Infratuzilma
      Config (pydantic)
      JWT + Redis blacklist
      Kafka publisher
      MinIO storage
      Distributed Lock
    Operatsiya
      Docker + entrypoint
      Alembic migration
      Health check
      Telegram xato xabari
```

## Patternlar ro'yxati
| Pattern | Qayerda |
| --- | --- |
| **Repository** | `repositories/` — DB so'rovlarini ajratadi |
| **Unit of Work** | `uow/` — transaction boshqaruvi |
| **Dependency Injection** | `dependencies/` + FastAPI `Depends` |
| **Port & Adapter (Hexagonal)** | `*/ports.py` (ABC) + `*/impl.py` / `clients/` |
| **DTO** | `schemas/` — model va API ni ajratadi |
| **Singleton** | `settings`, `redis_client`, `minio_client`, `publisher` |
| **Strategy / Adapter** | `MessagePublisher`, `FileStorage` |
| **Distributed Lock** | `core/redis_lock.py` |
| **Anti-Corruption Layer** | `clients/gateway.py` (tashqi JSON → ichki schema) |

---

# 19. To'liq so'rov hayot sikli (eng katta diagramma)

`GET /api/v1/uzassets/organizations/{tin}` so'rovi misolida:

```mermaid
sequenceDiagram
    autonumber
    participant C as Client
    participant MW as Middleware'lar
    participant R as organization router
    participant DI as get_org_service
    participant Uow as UnitOfWork
    participant Svc as OrganizationService
    participant Repo as OrganizationRepository
    participant DB as PostgreSQL

    C->>MW: GET /organizations/123456789
    MW->>MW: TrustedHost, CORS, Session, HTTPS
    MW->>R: route topiladi
    R->>DI: dependency'larni hal qiladi
    DI->>Uow: get_unit_of_work (session ochiladi)
    DI->>Svc: OrganizationService(uow, gateway)
    R->>Svc: get_org_by_tin("123456789")
    Svc->>Uow: async with uow
    Svc->>Repo: get_org_by_tin(tin)
    Repo->>DB: SELECT * FROM organizations WHERE tin=...
    DB-->>Repo: yozuv yoki None
    alt topilmadi
        Svc->>Svc: raise NotFoundException
        Svc-->>MW: app_exception_handler → 404 JSON
    else topildi
        Svc-->>R: Organization obyekti
        R->>R: OrganizationResponse (from_attributes)
        R-->>C: 200 + JSON
    end
    Note over Uow,DB: __aexit__: commit + session.close()
```

---

# Eslab qol (eng muhim 10 ta narsa)

1. **Qatlamli arxitektura:** Router → Service → UoW → Repository → Model. Har qatlam faqat pastdagini biladi.
2. **Endpointlar yupqa:** biznes logika service'da, endpoint faqat chaqiradi.
3. **Repository commit qilmaydi**, faqat `flush()`. Commit'ni Unit of Work boshqaradi.
4. **Unit of Work** bitta biznes amalini bitta transaction qiladi (commit yoki rollback).
5. **Dependency Injection** `Annotated[Type, Depends(...)]` alias bilan — kod toza va testlanadigan.
6. **Model ≠ Schema:** Model — DB jadvali, Schema — API formati (Pydantic).
7. **Abstract port + impl:** tashqi tizimlar (Kafka, MinIO, SSO, Gateway) interfeys ortida — almashtirsa bo'ladi.
8. **Maxsus exceptionlar** avtomatik to'g'ri HTTP status va JSON formatga aylanadi; production xatolari Telegram'ga ketadi.
9. **JWT + Redis blacklist** — logout qilingan token `jti` orqali bekor qilinadi.
10. **lifespan** startup/shutdown'ni boshqaradi; `/health` Kafka/Redis/Postgres'ni tekshiradi.

---

# Amaliyot (o'zing sinab ko'r)

1. **Yangi entity qo'sh:** `Region` modelini yarat (`models/region.py`), schema, repository, service, dependency va router yoz. `main.py` ga routerni qo'sh. Migratsiya yarat: `make make_migrations && make migrate`.
2. **UoW'ga repository qo'sh:** `uow/impl.py` ning `__aenter__` ida `self.region = RegionRepository(...)` qo'shib ko'r.
3. **Permission yarat:** `core/permissions.py` ga `get_admin_user` yozib, biror endpointga `AdminUserDep` qo'sh.
4. **Endpoint yoz:** `POST /organizations/` yaratish endpointini yoz — service ichida `uow.organization.add(...)` chaqir, TIN takrorlansa `BusinessException` ko'tar.
5. **Health check kengaytir:** `/health` ga MinIO tekshiruvini (hozir comment qilingan) qayta yoq.
6. **Kafka event yubor:** organization yaratilganda `kafka_publisher.publish(...)` chaqirib event yuborishni qo'sh.

---

*Bu hujjat `fastapi-structure/` loyihasining barcha fayllari to'liq o'qib chiqilib tayyorlandi.*
