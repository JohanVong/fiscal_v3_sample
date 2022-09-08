-- DROP SCHEMA dbo;;-- Drop table

-- DROP TABLE Address 

CREATE TABLE Address
(
    idAddress      integer primary key                           ,
    idTypeAddress  int                                          NULL,
    idTown         int                                          NULL,
    Street         nvarchar(200)  NULL,
    House          nvarchar(45)   NULL,
    Flat           nvarchar(45)   NULL,
    idTownExternal int                                          NULL

);;-- Drop table

-- DROP TABLE AgentRel 

CREATE TABLE AgentRel
(
    Id        integer primary key ,
    IdCompany int                NULL,
    IdAgents  int                NULL

);;-- Drop table

-- DROP TABLE Article 

CREATE TABLE Article
(
    idArticle      integer primary key                          ,
    Name           nvarchar(45)  NULL,
    QR             nvarchar(45)  NULL,
    idGroupArticle int                                         NULL,
    idCompany      int                                         NULL,
    idSection      int                                         NULL,
    Price          money                                       NULL,
    Discount       int                                         NULL,
    Markup         int                                         NULL,
    Active         bit                                         NULL
);;-- Drop table

-- DROP TABLE Balance 

CREATE TABLE Balance
(
    idBalance     integer primary key ,
    Amount        money              NULL,
    idTypeBalance int                NULL,
    idKKM         int                NULL
);-- Drop table

-- DROP TABLE Cashier 

CREATE TABLE Cashier
(
    idCashier       integer primary key                          ,
    idCompany       int                                         NULL,
    idUser          int                                         NULL,
    FIO             nvarchar(45)  NULL,
    idStatusCashier int                                         NULL,
    Lock            bit                                         NULL,
    idShift         int                                         NULL
);-- Drop table

-- DROP TABLE Company 

CREATE TABLE Company
(
    idCompany   integer primary key                           ,
    idUser      int                                          NULL,
    idOwnership int                                          NULL,
    BIN         nvarchar(45)   NULL,
    ShortName   nvarchar(256)  NULL,
    FullName    nvarchar(256)  NULL,
    FIO         nvarchar(256)  NULL,
    NDS         nvarchar(256)  NULL,
    idAddress   int                                          NULL
);-- Drop table

-- DROP TABLE Contact 

CREATE TABLE Contact
(
    idContact     integer primary key                          ,
    idTypeContact int                                         NULL,
    Name          nvarchar(45)  NULL
);-- Drop table

-- DROP TABLE Documents 

CREATE TABLE Documents
(
    idDocuments       integer primary key                          ,
    idShift           int                                         NULL,
    idUser            int                                         NULL,
    idTypeDocument    int                                         NULL,
    NumberDoc         varchar(255)  NULL,
    idDomain          int                                         NULL,
    DateDocument      datetime                              NULL,
    Cash              money                                       NULL,
    [Change]          money                                       NULL,
    Value             money                                       NULL,
    NonCash           money                                       NULL,
    Coins             int                                         NULL,
    FiscalNumber      bigint                                      NULL,
    AutonomousNumber  bigint                                      NULL,
    idKKM             int                                         NULL,
    CheckSum          varchar(255)  NULL,
    Offline           bit                                         NULL,
    DocChain          varchar(255)  NULL,
    CheckLink         varchar(255)  NULL,
    Uid               varchar(100)  NULL,
    idCompany         int                                         NULL,
    idDocumentsParent bigint                                      NULL,
    token             bigint                                      NULL,
    reqNum            int                                         NULL
);

CREATE INDEX [_dta_index_Documents_8_1589580701__K1] ON Documents (idDocuments);

CREATE INDEX [_dta_index_Documents_8_1589580701__K10] ON Documents (idDomain);

CREATE INDEX [_dta_index_Documents_8_1589580701__K10_K1] ON Documents (idDomain, idDocuments);

CREATE INDEX [_dta_index_Documents_8_1589580701__K10_K3_K2_K4] ON Documents (idDomain, idUser, idShift, idTypeDocument);

CREATE INDEX [_dta_index_Documents_8_1589580701__K1D_K2_K4] ON Documents (idDocuments DESC, idShift, idTypeDocument);

CREATE INDEX [_dta_index_Documents_8_1589580701__K1_17] ON Documents (idDocuments, FiscalNumber);

CREATE INDEX [_dta_index_Documents_8_1589580701__K1_K10] ON Documents (idDocuments, idDomain);

CREATE INDEX [_dta_index_Documents_8_1589580701__K1_K2] ON Documents (idDocuments, idShift);

CREATE INDEX [_dta_index_Documents_8_1589580701__K1_K3] ON Documents (idDocuments, idUser);

CREATE INDEX [_dta_index_Documents_8_1589580701__K1_K4] ON Documents (idDocuments, idTypeDocument);

CREATE INDEX [_dta_index_Documents_8_1589580701__K2] ON Documents (idShift);

CREATE INDEX [_dta_index_Documents_8_1589580701__K21] ON Documents (Offline);

CREATE INDEX [_dta_index_Documents_8_1589580701__K2_K1] ON Documents (idShift, idDocuments);

CREATE INDEX [_dta_index_Documents_8_1589580701__K2_K3_K4_K10] ON Documents (idShift, idUser, idTypeDocument, idDomain);

CREATE INDEX [_dta_index_Documents_8_1589580701__K2_K4] ON Documents (idShift, idTypeDocument);

CREATE INDEX [_dta_index_Documents_8_1589580701__K2_K4_K1D] ON Documents (idShift, idTypeDocument, idDocuments DESC);

CREATE INDEX [_dta_index_Documents_8_1589580701__K3] ON Documents (idUser);

CREATE INDEX [_dta_index_Documents_8_1589580701__K3_K1] ON Documents (idUser, idDocuments);

CREATE INDEX [_dta_index_Documents_8_1589580701__K3_K2_K4_K10] ON Documents (idUser, idShift, idTypeDocument, idDomain);

CREATE INDEX [_dta_index_Documents_8_1589580701__K4] ON Documents (idTypeDocument);

CREATE INDEX [_dta_index_Documents_8_1589580701__K4_K1] ON Documents (idTypeDocument, idDocuments);

CREATE INDEX [_dta_index_Documents_8_1589580701__K4_K3_K2_K10] ON Documents (idTypeDocument, idUser, idShift, idDomain);

CREATE INDEX [_dta_stat_1589580701_10_1_3_2] ON Documents (idDomain, idDocuments, idUser, idShift);

CREATE INDEX [_dta_stat_1589580701_10_3_2] ON Documents (idDomain, idUser, idShift);

CREATE INDEX [_dta_stat_1589580701_2_1] ON Documents (idShift, idDocuments);

CREATE INDEX [_dta_stat_1589580701_2_3_4_10] ON Documents (idShift, idUser, idTypeDocument, idDomain);

CREATE INDEX [_dta_stat_1589580701_3_2_4_10_1] ON Documents (idUser, idShift, idTypeDocument, idDomain, idDocuments);

CREATE INDEX [_dta_stat_1589580701_4_1_3] ON Documents (idTypeDocument, idDocuments, idUser);

CREATE INDEX [_dta_stat_1589580701_4_3] ON Documents (idTypeDocument, idUser);



-- DROP TABLE [Domain] 

CREATE TABLE [Domain]
(
    idDomain integer primary key                          ,
    Name     nvarchar(45)  NULL
);-- Drop table

-- DROP TABLE GroupArticle 

CREATE TABLE GroupArticle
(
    idGroupArticle integer primary key                          ,
    Name           nvarchar(45)  NULL,
    IDGroup        int                                         NULL
);-- Drop table

-- DROP TABLE GroupUser 

CREATE TABLE GroupUser
(
    idGroupUser integer primary key                          ,
    Name        nvarchar(45)  NULL
);-- Drop table

-- DROP TABLE KKM 

CREATE TABLE KKM
(
    idKKM        integer primary key                          ,
    idModeKKM    int                                         NULL,
    idOFD        int                                         NULL,
    idNK         int                                         NULL,
    idPlaceUsed  int                                         NULL,
    Name         nvarchar(45)  NULL,
    idStatusKKM  int                                         NULL,
    RNM          nvarchar(45)  NULL,
    Autorenew    bit                                         NULL,
    idCompany    int                                         NULL,
    idAddress    int                                         NULL,
    idCPCR       bigint                                      NULL,
    tokenCPCR    bigint                                      NULL,
    reqnumCPCR   int                                         NULL,
    OfflineQueue bigint                                      NOT NULL,
    ShiftIndex   int                                         NULL,
    Lock         bit                                         NOT NULL,
    IdSection    int                                         NULL,
    znm          int                                         NULL,
    idShift      int                                         NULL,
    IsActive     bit                                         NOT NULL
);




CREATE  INDEX ind_idSKKM ON KKM (idStatusKKM ASC)
    ;-- Drop table

-- DROP TABLE KKMParam 

CREATE TABLE KKMParam
(
    idKKMParam     integer primary key                          ,
    Value          nvarchar(45)  NULL,
    idTypeKKMParam int                                         NULL,
    idKKM          int                                         NULL,
    idTypeValue    int                                         NULL

);-- Drop table

-- DROP TABLE LKCasCon 

CREATE TABLE LKCasCon
(
    idLKCasCon integer primary key ,
    idContact  int                NULL,
    idCashier  int                NULL
);-- Drop table

-- DROP TABLE LKKKMAdd 

CREATE TABLE LKKKMAdd
(
    idLKKKMAdd integer primary key ,
    idAddress  int                NULL,
    idKKM      int                NULL
);-- Drop table

-- DROP TABLE LKMerAdd 

CREATE TABLE LKMerAdd
(
    idLKMerAdd integer primary key ,
    idCompany  int                NULL,
    idAddress  int                NULL
);-- Drop table

-- DROP TABLE LKMerCon 

CREATE TABLE LKMerCon
(
    idLKMerCon integer primary key ,
    idCompany  int                NULL,
    idContact  int                NULL);-- Drop table

-- DROP TABLE Migration 

CREATE TABLE Migration
(
    IDMigration       integer primary key,
    IDTerminal        int                NULL,
    IDMigrationStatus int                NULL);-- Drop table

-- DROP TABLE MigrationStatus 

CREATE TABLE MigrationStatus
(
    IDMigrationStatus int                                         NULL,
    StatusName        varchar(128)  NULL
);-- Drop table

-- DROP TABLE ModeKKM 

CREATE TABLE ModeKKM
(
    idModeKKM integer primary key        ,
    Name      nvarchar(45)  NULL);-- Drop table

-- DROP TABLE NK 

CREATE TABLE NK
(
    idNK     integer primary key         ,
    Name     nvarchar(45)  NULL,
    idRegion int                                         NULL,
    Code     int                                         NULL);-- Drop table

-- DROP TABLE OFD 

CREATE TABLE OFD
(
    idOFD integer primary key            ,
    Name  nvarchar(45)  NULL);-- Drop table

-- DROP TABLE Operations 

CREATE TABLE Operations
(
    idOperations  integer primary key,
    idDocuments   int                NULL,
    DateOperation datetime2(7)       NULL,
    Value         money              NULL);-- Drop table

-- DROP TABLE Ownership 

CREATE TABLE Ownership
(
    idOwnership integer primary key  ,
    Name        nvarchar(45)  NULL);-- Drop table

-- DROP TABLE ParamDoc 

CREATE TABLE ParamDoc
(
    idParamDoc     integer primary key,
    Value          nvarchar(45)  NULL,
    idTypeValue    int                                         NULL,
    idDocuments    int                                         NULL,
    idTypeParamDoc int                                         NULL);-- Drop table

-- DROP TABLE PermUser 

CREATE TABLE PermUser
(
    idPermUser   integer primary key  ,
    idGroupUser  int                                          NULL,
    idTypeObject int                                          NULL,
    NameObject   nvarchar(200)  NULL,
    idTypePerm   int                                          NULL);-- Drop table

-- DROP TABLE PlaceUsed 

CREATE TABLE PlaceUsed
(
    idPlaceUsed integer primary key                          ,
    Name        nvarchar(45)  NULL);-- Drop table

-- DROP TABLE [Position] 

CREATE TABLE [Position]
(
    idPosition  integer primary key                          ,
    idDocuments int                                          NULL,
    idArticle   int                                          NULL,
    [Number]    int                                          NULL,
    Price       money                                        NULL,
    Discount    money                                        NULL,
    Total       money                                        NULL,
    Markup      money                                        NULL,
    Qty         int                                          NULL,
    idSection   int                                          NULL,
    Nds         money                                        NULL,
    Storno      bit                                          NULL,
    NdsDiscount money                                        NULL,
    NdsMarkup   money                                        NULL,
    idCompany   int                                          NULL,
    Name        nvarchar(255)  NULL);-- Drop table

-- DROP TABLE Region 

CREATE TABLE Region
(
    idRegion integer primary key                          ,
    Name     nvarchar(45)  NULL);-- Drop table

-- DROP TABLE [Section] 

CREATE TABLE [Section]
(
    idSection integer primary key                         ,
    Name      nvarchar(45)  NULL,
    idKKM     int                                         NULL,
    NDS       int                                         NULL,
    Active    bit                                         NULL,
    idCompany int                                         NULL);-- Drop table

-- DROP TABLE Shift 

CREATE TABLE Shift
(
    idShift       integer primary key ,
    idUser        int                NULL,
    idKKM         int                NULL,
    idStatusShift int                NULL,
    DateOpen      datetime     NULL,
    DateClose     datetime     NULL,
    BalanceClose  money              NULL,
    BalanceOpen   money              NULL,
    ShiftIndex    int                NULL,
    idCompany     int                NULL);-- Drop table

-- DROP TABLE StatusCashier 

CREATE TABLE StatusCashier
(
    idStatusCashier integer primary key,
    Name            nvarchar(45)  NULL);-- Drop table

-- DROP TABLE StatusKKM 

CREATE TABLE StatusKKM
(
    idStatusKKM integer primary key                          ,
    Name        nvarchar(45)  NULL);-- Drop table

-- DROP TABLE StatusShift 

CREATE TABLE StatusShift
(
    idStatusShift integer primary key                        ,
    Name          nvarchar(45)  NULL);-- Drop table

-- DROP TABLE StatusUser 

CREATE TABLE StatusUser
(
    idStatusUser integer primary key                         ,
    Name         nvarchar(45)  NULL);-- Drop table

-- DROP TABLE Town 

CREATE TABLE Town
(
    idTown   integer primary key                          ,
    Name     nvarchar(45)  NULL,
    idRegion int                                         NULL,
    TimeZone int                                         NULL);-- Drop table

-- DROP TABLE TypeAddress 

CREATE TABLE TypeAddress
(
    idTypeAddress integer primary key                     ,
    Name          nvarchar(45)  NULL);-- Drop table

-- DROP TABLE TypeBalance 

CREATE TABLE TypeBalance
(
    idTypeBalance integer primary key                ,
    Name          nvarchar(45)  NULL);-- Drop table

-- DROP TABLE TypeContact 

CREATE TABLE TypeContact
(
    idTypeContact integer primary key               ,
    Name          nvarchar(45)  NULL);-- Drop table

-- DROP TABLE TypeDocument 

CREATE TABLE TypeDocument
(
    idTypeDocument integer primary key             ,
    Name           nvarchar(45)  NULL);-- Drop table

-- DROP TABLE TypeKKMParam 

CREATE TABLE TypeKKMParam
(
    idTypeKKMParam integer primary key   ,
    Name           nvarchar(45)  NULL);-- Drop table

-- DROP TABLE TypeObject 

CREATE TABLE TypeObject
(
    idTypeObject int                             ,
    Name         nvarchar(60)  NULL);-- Drop table

-- DROP TABLE TypeParamDoc 

CREATE TABLE TypeParamDoc
(
    idTypeParamDoc integer primary key                ,
    Name           nvarchar(45)  NULL,
    idTypeDocument int                                         NULL);-- Drop table

-- DROP TABLE TypePerm 

CREATE TABLE TypePerm
(
    idTypePerm integer primary key               ,
    Name       nvarchar(45)  NULL);-- Drop table

-- DROP TABLE TypeUser 

CREATE TABLE TypeUser
(
    idTypeUser integer primary key            ,
    Name       nvarchar(45)  NULL);-- Drop table

-- DROP TABLE TypeValue 

CREATE TABLE TypeValue
(
    idTypeValue integer primary key                  ,
    Name        nvarchar(45)  NULL);-- Drop table

-- DROP TABLE [User] 

CREATE TABLE [User]
(
    idUser       integer primary key                ,
    idTypeUser   int                                         NULL,
    PhoneLogin   nvarchar(45)  NULL,
    Password     nvarchar(90)  NULL,
    Name         nvarchar(45)  NULL,
    Lock         bit                                         NULL,
    idShift      int                                         NULL,
    idGroupUser  int                                         NULL,
    idStatusUser int                                         NULL);

CREATE UNIQUE INDEX UQ_User_PhoneLogin ON [User] (PhoneLogin);-- Drop table

-- DROP TABLE UserRel 

CREATE TABLE UserRel
(
    Id          integer primary key ,
    idUser      int                NULL,
    idCompany   int                NULL,
    idGroupUser int                NULL,
    Active      bit                NULL);-- Drop table

-- DROP TABLE ZReport 

CREATE TABLE ZReport
(
    id                      integer primary key,
    id_kkm                  int                NULL,
    date_open               datetime     NULL,
    date_close              datetime     NULL,
    balance_open            money              NULL,
    balance_close           money              NULL,
    count                   int                NULL,
    sales_qty               int                NULL,
    sales_amount            money              NULL,
    purchases_qty           int                NULL,
    purchases_amount        money              NULL,
    expenses_qty            int                NULL,
    expenses_amount         money              NULL,
    refunds_qty             int                NULL,
    refunds_amount          money              NULL,
    incomes_qty             int                NULL,
    incomes_amount          money              NULL,
    shift_index             int                NULL,
    idUser                  int                NULL,
    purchase_refunds_qty    int                NULL,
    purchase_refunds_amount money              NULL);-- Drop table


-- DROP TABLE system_config 

CREATE TABLE system_config
(
    id       integer primary key                         NOT NULL,
    [option] varchar(32)  NOT NULL,
    value    varchar(32)  NOT NULL)

