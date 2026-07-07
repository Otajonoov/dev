
# DDD asosida Identity & Access Management (IAM) Servisini Loyihalash

## 1. Strategic Design: Domainni tushunish

### 1.1 Problem Space va Solution Space

**Problem Space (Muammo maydoni):** Bizning tizimimiz foydalanuvchilar, ularning rollari va ruxsatlarini boshqarishi kerak. Bu **Identity and Access Management** domain'i.

**Solution Space (Yechim maydoni):** Biz ushbu muammoni hal qilish uchun aniq chegaralangan Bounded Context yaratamiz.

### 1.2 Bounded Context'ni aniqlash

```
┌─────────────────────────────────────────────────────────────┐
│          Identity & Access Management Context                │
│                                                               │
│  Ubiquitous Language:                                        │
│  - User (Foydalanuvchi)                                      │
│  - Role (Rol)                                                │
│  - Permission (Ruxsat)                                       │
│  - Tenant (Arендатор - multi-tenancy uchun)                 │
│  - Authentication (Autentifikatsiya)                         │
│  - Authorization (Avtorizatsiya)                             │
│  - Session (Sessiya)                                         │
└─────────────────────────────────────────────────────────────┘
```

**Nega Bounded Context kerak?**

- Har bir Bounded Context o'z Ubiquitous Language'ga ega
- "User" so'zi boshqa kontekstlarda (masalan, Billing Context) boshqa ma'noga ega bo'lishi mumkin
- Bu chegara bizning modelimizni boshqa tizimlardan izolyatsiya qiladi va murakkablikni boshqarishga yordam beradi

### 1.3 Subdomain'larni aniqlash

```
┌────────────────────────────────────────────────────────┐
│                    Domain Model                        │
├────────────────────────────────────────────────────────┤
│                                                        │
│  ┌──────────────────────────────────────┐            │
│  │   Core Domain (Asosiy yadro)         │            │
│  │   - Authorization Logic              │            │
│  │   - Role Management                  │            │
│  │   - Permission Management            │            │
│  └──────────────────────────────────────┘            │
│                                                        │
│  ┌──────────────────────────────────────┐            │
│  │   Supporting Subdomain               │            │
│  │   - User Registration                │            │
│  │   - Password Management              │            │
│  └──────────────────────────────────────┘            │
│                                                        │
│  ┌──────────────────────────────────────┐            │
│  │   Generic Subdomain                  │            │
│  │   - Token Generation (JWT)           │            │
│  │   - Email Notifications              │            │
│  └──────────────────────────────────────┘            │
└────────────────────────────────────────────────────────┘
```

**Tushuntirish:**

- **Core Domain**: Bu bizning raqobat ustunligimiz. Rollar va ruxsatlarni boshqarish mantiqimiz unique va biznesga maksimal qiymat beradi.
- **Supporting Subdomain**: Zarur, ammo tayyor yechim mavjud emas. O'zimiz yozishimiz kerak, lekin Core Domain qadar muhim emas.
- **Generic Subdomain**: JWT generatsiya qilish kabi umumiy funksionallik. Kutubxonalardan foydalanishimiz mumkin.

## 2. Tactical Design: Model Elementlari

### 2.1 Aggregate'larni aniqlash

**Nega Aggregate'lar kerak?** Aggregate'lar bizning transaksiya consistency chegaralarimizni belgilaydi. Har bir Aggregate bitta tranzaksiyada o'zgarishi va saqlanishi kerak bo'lgan elementlar klasteridir.

#### Aggregate #1: User Aggregate

go

```go
// Domain qatlami: domain/user/user.go

package user

import (
    "time"
    "errors"
)

// User - bu Aggregate Root
// Nega struct va class emas? Go'da class yo'q, struct ishlatamiz
type User struct {
    // Value Object - noyob identifikator
    id UserID
    
    // Value Object - tenant identifikatori
    tenantID TenantID
    
    // Value Object - email
    email Email
    
    // Value Object - password hash
    passwordHash PasswordHash
    
    // Entity - foydalanuvchi profili
    profile UserProfile
    
    // Value Object - status
    status UserStatus
    
    // Rol ID'lari - boshqa Aggregate'larga reference (faqat ID orqali!)
    roleIDs []RoleID
    
    // Audit fields
    createdAt time.Time
    updatedAt time.Time
    
    // Domain Event'lar - bu aggregate'da nima sodir bo'lganini qayd qiladi
    domainEvents []DomainEvent
}

// Value Object - UserID
// Nega alohida type? Primitiv obsession'dan qochish uchun
// string o'rniga UserID ishlatish xatolarni compile vaqtida topishga yordam beradi
type UserID struct {
    value string
}

func NewUserID(value string) (UserID, error) {
    if value == "" {
        return UserID{}, errors.New("user id cannot be empty")
    }
    return UserID{value: value}, nil
}

func (id UserID) String() string {
    return id.value
}

// Value Object - Email
// Nega Value Object? Chunki email'ning o'z validation mantiqi bor
type Email struct {
    value string
}

func NewEmail(value string) (Email, error) {
    // Email validation mantiqi - bu domain knowledge
    if !isValidEmail(value) {
        return Email{}, errors.New("invalid email format")
    }
    return Email{value: value}, nil
}

// Value Object - UserStatus
// Enum sifatida - faqat ma'lum qiymatlar
type UserStatus int

const (
    UserStatusPending UserStatus = iota
    UserStatusActive
    UserStatusSuspended
    UserStatusDeactivated
)

// Entity - UserProfile
// Nega Entity? Chunki u o'z identifikatsiyasiga ega va o'zgaruvchan
type UserProfile struct {
    firstName string
    lastName  string
    avatarURL string
}

// Factory Method - Aggregate yaratish uchun
// Nega constructor emas? Chunki biznes mantiqini inkapsulyatsiya qilish kerak
func CreateUser(
    tenantID TenantID,
    email Email,
    password string,
) (*User, error) {
    // Biznes qoidalarini tekshirish
    if password == "" {
        return nil, errors.New("password is required")
    }
    
    // Password hash qilish - bu ham domain knowledge
    passwordHash, err := NewPasswordHash(password)
    if err != nil {
        return nil, err
    }
    
    user := &User{
        id:           GenerateUserID(), // Domain service
        tenantID:     tenantID,
        email:        email,
        passwordHash: passwordHash,
        status:       UserStatusPending,
        roleIDs:      make([]RoleID, 0),
        createdAt:    time.Now(),
        domainEvents: make([]DomainEvent, 0),
    }
    
    // Domain Event yaratish
    // Nega kerak? Boshqa aggregate'lar yoki tizimlar bu haqida bilishi kerak
    user.addDomainEvent(UserRegisteredEvent{
        UserID:    user.id,
        TenantID:  user.tenantID,
        Email:     user.email,
        OccurredAt: time.Now(),
    })
    
    return user, nil
}

// Behavior Method - bu Aggregate'ning xatti-harakati
// Nega public setter emas? Anemic Domain Model'dan qochish uchun
// Barcha biznes mantiq aggregate ichida bo'lishi kerak
func (u *User) AssignRole(roleID RoleID) error {
    // Biznes qoidasini tekshirish
    if u.status != UserStatusActive {
        return errors.New("only active users can be assigned roles")
    }
    
    // Role allaqachon mavjud bo'lsa
    for _, existingRoleID := range u.roleIDs {
        if existingRoleID == roleID {
            return nil // Idempotent
        }
    }
    
    u.roleIDs = append(u.roleIDs, roleID)
    u.updatedAt = time.Now()
    
    // Domain Event
    u.addDomainEvent(UserRoleAssignedEvent{
        UserID:     u.id,
        RoleID:     roleID,
        OccurredAt: time.Now(),
    })
    
    return nil
}

func (u *User) Activate() error {
    // Biznes qoidasi: faqat pending user'larni activate qilish mumkin
    if u.status != UserStatusPending {
        return errors.New("only pending users can be activated")
    }
    
    u.status = UserStatusActive
    u.updatedAt = time.Now()
    
    u.addDomainEvent(UserActivatedEvent{
        UserID:     u.id,
        OccurredAt: time.Now(),
    })
    
    return nil
}

func (u *User) Authenticate(password string) error {
    if u.status != UserStatusActive {
        return errors.New("user is not active")
    }
    
    if !u.passwordHash.Matches(password) {
        // Domain Event - failed login attempt
        u.addDomainEvent(AuthenticationFailedEvent{
            UserID:     u.id,
            OccurredAt: time.Now(),
        })
        return errors.New("invalid credentials")
    }
    
    // Domain Event - successful authentication
    u.addDomainEvent(UserAuthenticatedEvent{
        UserID:     u.id,
        OccurredAt: time.Now(),
    })
    
    return nil
}

// Domain Event'larni boshqarish
func (u *User) addDomainEvent(event DomainEvent) {
    u.domainEvents = append(u.domainEvents, event)
}

func (u *User) DomainEvents() []DomainEvent {
    return u.domainEvents
}

func (u *User) ClearDomainEvents() {
    u.domainEvents = make([]DomainEvent, 0)
}
```

**Muhim tushuntirishlar:**

1. **Nega User Aggregate Root?**
    - U o'z lifecycle'ni boshqaradi
    - U transaksiya consistency chegarasini belgilaydi
    - Barcha o'zgarishlar User orqali o'tadi
2. **Nega roleIDs faqat ID'lar?**
    - Aggregate'larni kichik saqlash uchun (Rule #2)
    - Boshqa aggregate'larga faqat identifikator orqali murojaat (Rule #3)
    - Bu xotirada Role aggregate'ni yuklamasdan ishlash imkonini beradi
3. **Nega private field'lar?**
    - Encapsulation - biznes mantiqni himoya qilish
    - To'g'ridan-to'g'ri field'larga kirish Anemic Domain Model yaratadi

#### Aggregate #2: Role Aggregate

go

```go
// domain/role/role.go

package role

import (
    "time"
    "errors"
)

// Role Aggregate Root
type Role struct {
    id          RoleID
    tenantID    TenantID
    name        RoleName
    description string
    
    // Permission'larga reference - faqat ID'lar
    permissionIDs []PermissionID
    
    // Metadata
    isSystemRole bool // System role'larni o'chirish mumkin emas
    createdAt    time.Time
    updatedAt    time.Time
    
    domainEvents []DomainEvent
}

type RoleName struct {
    value string
}

func NewRoleName(value string) (RoleName, error) {
    if value == "" {
        return RoleName{}, errors.New("role name cannot be empty")
    }
    // System reserved nomlarni tekshirish
    if isReservedRoleName(value) {
        return RoleName{}, errors.New("role name is reserved")
    }
    return RoleName{value: value}, nil
}

// Factory Method
func CreateRole(
    tenantID TenantID,
    name RoleName,
    description string,
) (*Role, error) {
    role := &Role{
        id:            GenerateRoleID(),
        tenantID:      tenantID,
        name:          name,
        description:   description,
        permissionIDs: make([]PermissionID, 0),
        isSystemRole:  false,
        createdAt:     time.Now(),
        domainEvents:  make([]DomainEvent, 0),
    }
    
    role.addDomainEvent(RoleCreatedEvent{
        RoleID:     role.id,
        TenantID:   role.tenantID,
        Name:       role.name,
        OccurredAt: time.Now(),
    })
    
    return role, nil
}

// Behavior - permission qo'shish
func (r *Role) GrantPermission(permissionID PermissionID) error {
    // Biznes qoidasi: System role'larni o'zgartirish mumkin emas
    if r.isSystemRole {
        return errors.New("cannot modify system roles")
    }
    
    // Idempotent - allaqachon mavjud bo'lsa
    for _, existingPermID := range r.permissionIDs {
        if existingPermID == permissionID {
            return nil
        }
    }
    
    r.permissionIDs = append(r.permissionIDs, permissionID)
    r.updatedAt = time.Now()
    
    r.addDomainEvent(PermissionGrantedToRoleEvent{
        RoleID:       r.id,
        PermissionID: permissionID,
        OccurredAt:   time.Now(),
    })
    
    return nil
}

func (r *Role) RevokePermission(permissionID PermissionID) error {
    if r.isSystemRole {
        return errors.New("cannot modify system roles")
    }
    
    // Permission'ni topish va o'chirish
    for i, existingPermID := range r.permissionIDs {
        if existingPermID == permissionID {
            // Slice'dan o'chirish
            r.permissionIDs = append(
                r.permissionIDs[:i], 
                r.permissionIDs[i+1:]...,
            )
            r.updatedAt = time.Now()
            
            r.addDomainEvent(PermissionRevokedFromRoleEvent{
                RoleID:       r.id,
                PermissionID: permissionID,
                OccurredAt:   time.Now(),
            })
            
            return nil
        }
    }
    
    return errors.New("permission not found in role")
}

// Query method - permission borligini tekshirish
func (r *Role) HasPermission(permissionID PermissionID) bool {
    for _, existingPermID := range r.permissionIDs {
        if existingPermID == permissionID {
            return true
        }
    }
    return false
}
```

#### Aggregate #3: Permission Aggregate

go

```go
// domain/permission/permission.go

package permission

import "time"

// Permission Aggregate Root
// Bu aggregate juda oddiy - faqat ma'lumot saqlaydi
type Permission struct {
    id          PermissionID
    resource    Resource      // masalan: "user", "role", "document"
    action      Action        // masalan: "create", "read", "update", "delete"
    description string
    createdAt   time.Time
}

// Value Object - Resource
type Resource struct {
    value string
}

// Value Object - Action  
type Action struct {
    value string
}

const (
    ActionCreate = "create"
    ActionRead   = "read"
    ActionUpdate = "update"
    ActionDelete = "delete"
)

func CreatePermission(
    resource Resource,
    action Action,
    description string,
) *Permission {
    return &Permission{
        id:          GeneratePermissionID(),
        resource:    resource,
        action:      action,
        description: description,
        createdAt:   time.Now(),
    }
}

// Query method
func (p *Permission) Matches(resource Resource, action Action) bool {
    return p.resource == resource && p.action == action
}
```

**Nega Permission aggregate bunday oddiy?**

- Barcha aggregate'lar murakkab bo'lishi shart emas
- Permission faqat reference data
- Uning asosiy vazifasi - identity va immutability

### 2.2 Domain Event'lar

go

```go
// domain/events/user_events.go

package events

import "time"

// DomainEvent interface - barcha event'lar uchun
type DomainEvent interface {
    OccurredOn() time.Time
    EventType() string
}

// UserRegisteredEvent - foydalanuvchi ro'yxatdan o'tdi
type UserRegisteredEvent struct {
    UserID     UserID
    TenantID   TenantID
    Email      Email
    OccurredAt time.Time
}

func (e UserRegisteredEvent) OccurredOn() time.Time {
    return e.OccurredAt
}

func (e UserRegisteredEvent) EventType() string {
    return "UserRegistered"
}

// UserActivatedEvent - foydalanuvchi aktivlashtirildi
type UserActivatedEvent struct {
    UserID     UserID
    OccurredAt time.Time
}

func (e UserActivatedEvent) OccurredOn() time.Time {
    return e.OccurredAt
}

func (e UserActivatedEvent) EventType() string {
    return "UserActivated"
}

// UserRoleAssignedEvent - foydalanuvchiga rol tayinlandi
type UserRoleAssignedEvent struct {
    UserID     UserID
    RoleID     RoleID
    OccurredAt time.Time
}

func (e UserRoleAssignedEvent) OccurredOn() time.Time {
    return e.OccurredAt
}

func (e UserRoleAssignedEvent) EventType() string {
    return "UserRoleAssigned"
}

// UserAuthenticatedEvent - foydalanuvchi autentifikatsiya qildi
type UserAuthenticatedEvent struct {
    UserID     UserID
    OccurredAt time.Time
}

func (e UserAuthenticatedEvent) OccurredOn() time.Time {
    return e.OccurredAt
}

func (e UserAuthenticatedEvent) EventType() string {
    return "UserAuthenticated"
}
```

**Domain Event'lar nima uchun kerak?**

1. **Eventual Consistency uchun**: Bir aggregate o'zgarganda, boshqa aggregate'lar buni bilishi kerak
2. **Audit Trail**: Sistema da nima sodir bo'lganini yozib borish
3. **Integration**: Boshqa Bounded Context'lar bu event'larni subscribe qilishi mumkin
4. **Business Intelligence**: Event'lar asosida analytics qilish mumkin

### 2.3 Domain Service'lar

go

```go
// domain/services/authorization_service.go

package services

// AuthorizationService - Domain Service
// Nega Domain Service? Chunki bu mantiq bitta aggregate'ga tegishli emas
// Bu bir nechta aggregate'lar o'rtasidagi munosabatlarni boshqaradi
type AuthorizationService interface {
    // User'ning ma'lum resource'ga action qilish huquqi bormi?
    IsAuthorized(
        userID UserID, 
        resource Resource, 
        action Action,
    ) (bool, error)
}

type authorizationService struct {
    userRepo       UserRepository
    roleRepo       RoleRepository
    permissionRepo PermissionRepository
}

func NewAuthorizationService(
    userRepo UserRepository,
    roleRepo RoleRepository,
    permissionRepo PermissionRepository,
) AuthorizationService {
    return &authorizationService{
        userRepo:       userRepo,
        roleRepo:       roleRepo,
        permissionRepo: permissionRepo,
    }
}

func (s *authorizationService) IsAuthorized(
    userID UserID,
    resource Resource,
    action Action,
) (bool, error) {
    // 1. User'ni yuklash
    user, err := s.userRepo.FindByID(userID)
    if err != nil {
        return false, err
    }
    
    // 2. User active emasligini tekshirish
    if user.Status() != UserStatusActive {
        return false, nil
    }
    
    // 3. User'ning barcha role'larini yuklash
    roles := make([]*Role, 0)
    for _, roleID := range user.RoleIDs() {
        role, err := s.roleRepo.FindByID(roleID)
        if err != nil {
            continue // Bu role topilmasa, davom et
        }
        roles = append(roles, role)
    }
    
    // 4. Har bir role'ning permission'larini tekshirish
    for _, role := range roles {
        for _, permID := range role.PermissionIDs() {
            permission, err := s.permissionRepo.FindByID(permID)
            if err != nil {
                continue
            }
            
            // Permission mos keladimi?
            if permission.Matches(resource, action) {
                return true, nil
            }
        }
    }
    
    return false, nil
}
```

**Nega Domain Service?**

- Bu mantiq faqat User, Role yoki Permission aggregate'ga tegishli emas
- Bu uchala aggregate'ni birga ishlatadi
- Bu domain knowledge, ammo bitta aggregate'ga sig'maydi

### 2.4 Repository Interface'lari (Domain qatlamida)

go

```go
// domain/repositories/user_repository.go

package repositories

// UserRepository - bu interface, implementation emas!
// Nega? Repository pattern - domain va infrastructure'ni ajratish uchun
type UserRepository interface {
    // Aggregate'ni ID orqali topish
    FindByID(id UserID) (*User, error)
    
    // Email orqali topish (unique constraint)
    FindByEmail(tenantID TenantID, email Email) (*User, error)
    
    // Aggregate'ni saqlash
    Save(user *User) error
    
    // Aggregate'ni o'chirish
    Delete(id UserID) error
    
    // Tenant bo'yicha barcha user'lar
    FindByTenant(tenantID TenantID) ([]*User, error)
}

// RoleRepository
type RoleRepository interface {
    FindByID(id RoleID) (*Role, error)
    FindByName(tenantID TenantID, name RoleName) (*Role, error)
    Save(role *Role) error
    Delete(id RoleID) error
    FindByTenant(tenantID TenantID) ([]*Role, error)
}

// PermissionRepository  
type PermissionRepository interface {
    FindByID(id PermissionID) (*Permission, error)
    FindByResourceAndAction(resource Resource, action Action) (*Permission, error)
    Save(permission *Permission) error
    FindAll() ([]*Permission, error)
}
```

**Repository Pattern nima uchun?**

- Domain qatlami ma'lumotlar bazasi haqida bilmasligi kerak
- Interface domain'da, implementation infrastructure'da
- Test qilishda mock repository'lardan foydalanish oson
- Ma'lumotlar bazasini o'zgartirish domain'ni o'zgartirmaydi

## 3. Application Layer (Ilova qatlami)

go

```go
// application/commands/register_user_command.go

package commands

// Command - bu foydalanuvchining niyatini ifodalaydi
// Nega Command pattern? CQRS ning Command qismi
type RegisterUserCommand struct {
    TenantID string
    Email    string
    Password string
}

// CommandHandler - command'ni bajaradi
type RegisterUserCommandHandler struct {
    userRepo      UserRepository
    eventPublisher EventPublisher // Event'larni nashr qilish uchun
}

func NewRegisterUserCommandHandler(
    userRepo UserRepository,
    eventPublisher EventPublisher,
) *RegisterUserCommandHandler {
    return &RegisterUserCommandHandler{
        userRepo:      userRepo,
        eventPublisher: eventPublisher,
    }
}

func (h *RegisterUserCommandHandler) Handle(
    cmd RegisterUserCommand,
) error {
    // 1. Value Object'larni yaratish va validate qilish
    tenantID, err := NewTenantID(cmd.TenantID)
    if err != nil {
        return err
    }
    
    email, err := NewEmail(cmd.Email)
    if err != nil {
        return err
    }
    
    // 2. Email allaqachon mavjudligini tekshirish
    existingUser, _ := h.userRepo.FindByEmail(tenantID, email)
    if existingUser != nil {
        return errors.New("email already registered")
    }
    
    // 3. Aggregate yaratish - bu yerda biznes mantiq ishlaydi
    user, err := CreateUser(tenantID, email, cmd.Password)
    if err != nil {
        return err
    }
    
    // 4. Aggregate'ni saqlash
    if err := h.userRepo.Save(user); err != nil {
        return err
    }
    
    // 5. Domain Event'larni nashr qilish
    // Bu eventual consistency uchun muhim!
    for _, event := range user.DomainEvents() {
        h.eventPublisher.Publish(event)
    }
    user.ClearDomainEvents()
    
    return nil
}
```

go

```go
// application/commands/assign_role_command.go

type AssignRoleToUserCommand struct {
    UserID string
    RoleID string
}

type AssignRoleToUserCommandHandler struct {
    userRepo UserRepository
    roleRepo RoleRepository
    eventPublisher EventPublisher
}

func (h *AssignRoleToUserCommandHandler) Handle(
    cmd AssignRoleToUserCommand,
) error {
    // 1. Value Object'lar
    userID, err := NewUserID(cmd.UserID)
    if err != nil {
        return err
    }
    
    roleID, err := NewRoleID(cmd.RoleID)
    if err != nil {
        return err
    }
    
    // 2. User aggregate'ni yuklash
    user, err := h.userRepo.FindByID(userID)
    if err != nil {
        return err
    }
    
    // 3. Role mavjudligini tekshirish
    role, err := h.roleRepo.FindByID(roleID)
    if err != nil {
        return err
    }
    
    // 4. Biznes mantiqni bajarish (aggregate method)
    if err := user.AssignRole(role.ID()); err != nil {
        return err
    }
    
    // 5. Saqlash
    if err := h.userRepo.Save(user); err != nil {
        return err
    }
    
    // 6. Event'larni nashr qilish
    for _, event := range user.DomainEvents() {
        h.eventPublisher.Publish(event)
    }
    user.ClearDomainEvents()
    
    return nil
}
```

go

```go
// application/queries/get_user_permissions_query.go

package queries

// Query - o'qish operatsiyasi (CQRS ning Query qismi)
type GetUserPermissionsQuery struct {
    UserID string
}

// Query Result DTO (Data Transfer Object)
type UserPermissionsDTO struct {
    UserID      string
    Email       string
    Roles       []RoleDTO
    Permissions []PermissionDTO
}

type RoleDTO struct {
    RoleID      string
    Name        string
    Description string
}

type PermissionDTO struct {
    PermissionID string
    Resource     string
    Action       string
    Description  string
}

// QueryHandler
type GetUserPermissionsQueryHandler struct {
    userRepo       UserRepository
    roleRepo       RoleRepository
    permissionRepo PermissionRepository
}

func (h *GetUserPermissionsQueryHandler) Handle(
    query GetUserPermissionsQuery,
) (*UserPermissionsDTO, error) {
    // 1. User'ni yuklash
    userID, _ := NewUserID(query.UserID)
    user, err := h.userRepo.FindByID(userID)
    if err != nil {
        return nil, err
    }
    
    // 2. Role'larni yuklash
    roles := make([]RoleDTO, 0)
    allPermissions := make(map[PermissionID]*Permission)
    
    for _, roleID := range user.RoleIDs() {
        role, err := h.roleRepo.FindByID(roleID)
        if err != nil {
            continue
        }
        
        roles = append(roles, RoleDTO{
            RoleID:      role.ID().String(),
            Name:        role.Name().String(),
            Description: role.Description(),
        })
        
        // 3. Har bir role'ning permission'larini yig'ish
        for _, permID := range role.PermissionIDs() {
            perm, err := h.permissionRepo.FindByID(permID)
            if err != nil {
                continue
            }
            allPermissions[permID] = perm
        }
    }
    
    // 4. Permission'larni DTO'ga o'zgartirish
    permissions := make([]PermissionDTO, 0, len(allPermissions))
    for _, perm := range allPermissions {
        permissions = append(permissions, PermissionDTO{
            PermissionID: perm.ID().String(),
            Resource:     perm.Resource().String(),
            Action:       perm.Action().String(),
            Description:  perm.Description(),
        })
    }
    
    // 5. DTO'ni qaytarish
    return &UserPermissionsDTO{
        UserID:      user.ID().String(),
        Email:       user.Email().String(),
        Roles:       roles,
        Permissions: permissions,
    }, nil
}
```

**Application Layer nima uchun kerak?**

1. **Use Case Orchestration**: Bir use case bir nechta aggregate'lar bilan ishlashi mumkin
2. **Transaction Management**: Application layer tranzaksiyalarni boshqaradi
3. **Event Publishing**: Domain event'larni messaging tizimiga yuboradi
4. **DTO Transformation**: Domain modellarni API uchun DTO'larga aylantiradi
5. **CQRS**: Command va Query'larni ajratadi

## 4. Infrastructure Layer

go

```go
// infrastructure/persistence/postgres/user_repository_impl.go

package postgres

import (
    "database/sql"
)

// PostgresUserRepository - UserRepository interface'ni amalga oshiradi
type PostgresUserRepository struct {
    db *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) *PostgresUserRepository {
    return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) FindByID(id UserID) (*User, error) {
    query := `
        SELECT id, tenant_id, email, password_hash, status, 
               first_name, last_name, created_at, updated_at
        FROM users
        WHERE id = $1
    `
    
    var (
        userID       string
        tenantID     string
        email        string
        passwordHash string
        status       int
        firstName    string
        lastName     string
        createdAt    time.Time
        updatedAt    time.Time
    )
    
    err := r.db.QueryRow(query, id.String()).Scan(
        &userID, &tenantID, &email, &passwordHash, &status,
        &firstName, &lastName, &createdAt, &updatedAt,
    )
    if err != nil {
        return nil, err
    }
    
    // Database qatorini Domain model'ga aylantirish (Reconstitution)
    user := &User{
        id:       NewUserIDFromString(userID),
        tenantID: NewTenantIDFromString(tenantID),
        email:    NewEmailFromString(email),
        passwordHash: NewPasswordHashFromString(passwordHash),
        status:   UserStatus(status),
        profile: UserProfile{
            firstName: firstName,
            lastName:  lastName,
        },
        createdAt: createdAt,
        updatedAt: updatedAt,
        roleIDs:   []RoleID{}, // Alohida jadvaldan yuklash kerak
        domainEvents: []DomainEvent{},
    }
    
    // User role'larini yuklash (join table)
    roleQuery := `
        SELECT role_id 
        FROM user_roles 
        WHERE user_id = $1
    `
    rows, err := r.db.Query(roleQuery, userID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    for rows.Next() {
        var roleIDStr string
        if err := rows.Scan(&roleIDStr); err != nil {
            continue
        }
        user.roleIDs = append(user.roleIDs, NewRoleIDFromString(roleIDStr))
    }
    
    return user, nil
}

func (r *PostgresUserRepository) Save(user *User) error {
    // Transaction boshlanishi
    tx, err := r.db.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()
    
    // User'ni saqlash yoki yangilash
    query := `
        INSERT INTO users (
            id, tenant_id, email, password_hash, status,
            first_name, last_name, created_at, updated_at
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
        ON CONFLICT (id) DO UPDATE SET
            email = EXCLUDED.email,
            password_hash = EXCLUDED.password_hash,
            status = EXCLUDED.status,
            first_name = EXCLUDED.first_name,
            last_name = EXCLUDED.last_name,
            updated_at = EXCLUDED.updated_at
    `
    
    _, err = tx.Exec(query,
        user.ID().String(),
        user.TenantID().String(),
        user.Email().String(),
        user.PasswordHash().String(),
        int(user.Status()),
        user.Profile().FirstName(),
        user.Profile().LastName(),
        user.CreatedAt(),
        user.UpdatedAt(),
    )
    if err != nil {
        return err
    }
    
    // User role'larini saqlash
    // Avval eski role'larni o'chirish
    deleteRolesQuery := `DELETE FROM user_roles WHERE user_id = $1`
    _, err = tx.Exec(deleteRolesQuery, user.ID().String())
    if err != nil {
        return err
    }
    
    // Yangi role'larni qo'shish
    insertRoleQuery := `
        INSERT INTO user_roles (user_id, role_id) 
        VALUES ($1, $2)
    `
    for _, roleID := range user.RoleIDs() {
        _, err = tx.Exec(insertRoleQuery, 
            user.ID().String(), 
            roleID.String(),
        )
        if err != nil {
            return err
        }
    }
    
    // Transaction commit
    return tx.Commit()
}
```

**Infrastructure Layer nima uchun kerak?**

- Domain qatlami texnik detallar haqida bilmasligi kerak
- Ma'lumotlar bazasi, messaging, external API'lar - barchasi shu yerda
- Implementation'ni o'zgartirish domain'ni buzмайdi

## 5. Interface Layer (API/Presentation)

go

````go
// interfaces/http/handlers/user_handler.go

package handlers

import (
    "encoding/json"
    "net/http"
)

type UserHandler struct {
    registerUserHandler *RegisterUserCommandHandler
    getUserPermHandler  *GetUserPermissionsQueryHandler
}

// HTTP Request DTO
type RegisterUserRequest struct {
    TenantID string `json:"tenant_id"`
    Email    string `json:"email"`
    Password string `json:"password"`
}

// HTTP Response DTO
type RegisterUserResponse struct {
    UserID  string `json:"user_id"`
    Message string `json:"message"`
}

func (h *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
    // 1. Request'ni parse qilish
    var req RegisterUserRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }
    
    // 2. Command yaratish
    cmd := RegisterUserCommand{
        TenantID: req.TenantID,
        Email:    req.Email,
        Password: req.Password,
    }
    
    // 3. Command'ni bajarish
    err := h.registerUserHandler.Handle(cmd)
    if err != nil {
        // Biznes xatolarini HTTP status code'larga map qilish
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    // 4. Response qaytarish
    resp := RegisterUserResponse{
        Message: "User registered successfully",
    }
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteStatus(http.StatusCreated)
    json.NewEncoder(w).Encode(resp)
}
```

## 6. Arxitektura Qatlamlari - To'liq Rasm
```
┌──────────────────────────────────────────────────────────┐
│                    Interface Layer                        │
│  (HTTP Handlers, gRPC Services, CLI Commands)            │
│  - Request/Response DTO'lar                              │
│  - HTTP status code mapping                               │
│  - Authentication middleware                              │
└────────────────────┬─────────────────────────────────────┘
                     │
                     ▼
┌──────────────────────────────────────────────────────────┐
│                  Application Layer                        │
│  - Command Handlers (write operations)                   │
│  - Query Handlers (read operations)                      │
│  - Use Case orchestration                                │
│  - Transaction management                                │
│  - Event publishing                                      │
└────────────────────┬─────────────────────────────────────┘
                     │
                     ▼
┌──────────────────────────────────────────────────────────┐
│                    Domain Layer                           │
│  ┌────────────────────────────────────────────────────┐ │
│  │  Aggregates (User, Role, Permission)               │ │
│  │  - Aggregate Roots                                 │ │
│  │  - Entities                                        │ │
│  │  - Value Objects                                   │ │
│  │  - Business Rules                                  │ │
│  │  - Domain Events                                   │ │
│  └────────────────────────────────────────────────────┘ │
│  ┌────────────────────────────────────────────────────┐ │
│  │  Domain Services                                   │ │
│  │  - Authorization Service                           │ │
│  └────────────────────────────────────────────────────┘ │
│  ┌────────────────────────────────────────────────────┐ │
│  │  Repository Interfaces                             │ │
│  │  - UserRepository                                  │ │
│  │  - RoleRepository                                  │ │
│  │  - PermissionRepository                            │ │
│  └────────────────────────────────────────────────────┘ │
└────────────────────┬─────────────────────────────────────┘
                     │
                     ▼
┌──────────────────────────────────────────────────────────┐
│                Infrastructure Layer                       │
│  - Repository Implementations (PostgreSQL, MongoDB)      │
│  - Event Publisher (RabbitMQ, Kafka)                     │
│  - External Services (Email, SMS)                        │
│  - Caching (Redis)                                       │
└──────────────────────────────────────────────────────────┘
```

## 7. Context Map - Boshqa tizimlar bilan integratsiya
```
┌─────────────────────────────────────────────────────────┐
│         Identity & Access Management Context             │
│                  (Bizning tizimimiz)                     │
└────────────────────┬────────────────────────────────────┘
                     │
                     │ Domain Event'lar (User Registered, etc)
                     │ Published Language (JSON Schema)
                     │
          ┌──────────┴──────────┬────────────────┐
          │                     │                 │
          ▼                     ▼                 ▼
┌──────────────────┐  ┌─────────────────┐  ┌────────────┐
│ Notification     │  │ Audit Log       │  │ Analytics  │
│ Context          │  │ Context         │  │ Context    │
│                  │  │                 │  │            │
│ (Conformist)     │  │ (Conformist)    │  │(Conformist)│
└──────────────────┘  └─────────────────┘  └────────────┘

┌──────────────────────────────────────────────────────────┐
│  Bizning Context - Open Host Service sifatida            │
│  - REST API endpoints                                     │
│  - Published Language: JSON                               │
│  - Anticorruption Layer kerak emas (biz upstream'miz)    │
└──────────────────────────────────────────────────────────┘
```

## 8. Asosiy DDD Printsiplarining Qo'llanishi

### 8.1 Ubiquitous Language
```
✓ User, Role, Permission - domain ekspertlari ishlatadigan atamalar
✓ "Assign Role" - "add role to user" emas
✓ "Grant Permission" - "insert permission" emas
✓ Kod domain tilida yozilgan
````

### 8.2 Aggregate Rules

**Rule #1: Biznes invariantlarni himoya qiling**

go

```go
// ✓ To'g'ri: Aggregate ichida tekshirish
func (u *User) Activate() error {
    if u.status != UserStatusPending {
        return errors.New("only pending users can be activated")
    }
    // ...
}

// ✗ Noto'g'ri: Tekshirish tashqarida
func ActivateUser(user *User) {
    user.SetStatus(UserStatusActive) // Hech qanday tekshirish yo'q!
}
```

**Rule #2: Kichik aggregate'lar**

go

```go
// ✓ To'g'ri: User aggregate'da faqat roleID'lar
type User struct {
    roleIDs []RoleID // Faqat ID'lar
}

// ✗ Noto'g'ri: Butun Role aggregate'larni yuklash
type User struct {
    roles []*Role // Katta aggregate, sekin yuklash
}
```

**Rule #3: ID orqali murojaat**

go

```go
// ✓ To'g'ri
user.AssignRole(roleID) // Faqat ID kerak

// ✗ Noto'g'ri
user.AssignRole(role) // Butun aggregate kerak emas
```

**Rule #4: Eventual Consistency**

go

````go
// Domain Event orqali boshqa aggregate'ni yangilash
user.AssignRole(roleID)
// Event: UserRoleAssignedEvent nashr etiladi
// Bu event role statistikasini yangilashi mumkin (alohida tranzaksiyada)
```

### 8.3 Layered Architecture - Dependency Rule
```
Interface Layer  ──────┐
                       │
Application Layer ─────┤ Faqat domain'ga bog'liq
                       │
Domain Layer     ──────┘ Hech kimga bog'liq emas!
                       
Infrastructure Layer ───> Domain'ni implement qiladi
````

**Nega bu muhim?**

- Domain layer business logic - eng qimmatli qism
- U texnik detallardan izolyatsiya qilingan
- Test qilish oson
- Texnologiyani o'zgartirish domain'ni buzмайди

## Xulosa

Bu loyiha DDD ning barcha asosiy printsiрларini namoyish etadi:

1. **Strategic Design**:
    - Bounded Context aniq belgilangan
    - Ubiquitous Language'dan foydalanilgan
    - Subdomain'lar ajratilgan (Core, Supporting, Generic)
2. **Tactical Design**:
    - Aggregate'lar to'g'ri loyihalangan (kichik, consistency chegaralari aniq)
    - Entity va Value Object'lar farqlanadi
    - Domain Event'lar eventual consistency uchun ishlatiladi
    - Repository pattern domain'ni infrastructure'dan ajratadi
3. **Architecture**:
    - Layered Architecture (Dependency Rule)
    - CQRS (Command va Query ajratish)
    - Domain qatlami markazda va mustaqil
4. **Best Practices**:
    - Anemic Domain Model'dan qochish
    - Encapsulation (private field'lar, behavior method'lar)
    - Immutable Value Object'lar
    - Factory Method'lar aggregate yaratish uchun

Bu yondashuv biznesga qiymat beradi, chunki:

- Kod biznes tilida yozilgan
- O'zgarishlar oson kiritiladi
- Texnologiya o'zgarishi qo'rqinchli emas
- Test qilish oddiy
- Yangi developer'lar tezda tushunadi (Ubiquitous Language tufayli)
