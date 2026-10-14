package nex

import (
	"fmt"
	"os"
	"strconv"

	nex "github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	common_datastore "github.com/PretendoNetwork/nex-protocols-common-go/v2/datastore"
	common_globals "github.com/PretendoNetwork/nex-protocols-common-go/v2/globals"
	common_matchmaking "github.com/PretendoNetwork/nex-protocols-common-go/v2/match-making"
	common_matchmaking_ext "github.com/PretendoNetwork/nex-protocols-common-go/v2/match-making-ext"
	common_matchmake_extension "github.com/PretendoNetwork/nex-protocols-common-go/v2/matchmake-extension"
	common_nat_traversal "github.com/PretendoNetwork/nex-protocols-common-go/v2/nat-traversal"
	common_ranking "github.com/PretendoNetwork/nex-protocols-common-go/v2/ranking"
	common_secure "github.com/PretendoNetwork/nex-protocols-common-go/v2/secure-connection"
	common_utility "github.com/PretendoNetwork/nex-protocols-common-go/v2/utility"
	datastore "github.com/PretendoNetwork/nex-protocols-go/v2/datastore"
	matchmaking "github.com/PretendoNetwork/nex-protocols-go/v2/match-making"
	matchmaking_ext "github.com/PretendoNetwork/nex-protocols-go/v2/match-making-ext"
	matchmaking_types "github.com/PretendoNetwork/nex-protocols-go/v2/match-making/types"
	matchmake_extension "github.com/PretendoNetwork/nex-protocols-go/v2/matchmake-extension"
	nat_traversal "github.com/PretendoNetwork/nex-protocols-go/v2/nat-traversal"
	ranking "github.com/PretendoNetwork/nex-protocols-go/v2/ranking"
	secure "github.com/PretendoNetwork/nex-protocols-go/v2/secure-connection"
	utility "github.com/PretendoNetwork/nex-protocols-go/v2/utility"
	"github.com/Protarium-Network/mighty-no9-nex/database"
	"github.com/Protarium-Network/mighty-no9-nex/globals"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var SecureServer *nex.PRUDPServer
var SecureEndpoint *nex.PRUDPEndPoint

func StartSecureServer() {
	SecureServer = nex.NewPRUDPServer()

	// See authentication.go.
	SecureServer.PRUDPV1Settings.LegacyConnectionSignature = false

	SecureEndpoint = nex.NewPRUDPEndPoint(1)
	SecureEndpoint.IsSecureEndPoint = true
	SecureEndpoint.ServerAccount = globals.SecureServerAccount
	SecureEndpoint.AccountDetailsByPID = globals.AccountDetailsByPID
	SecureEndpoint.AccountDetailsByUsername = globals.AccountDetailsByUsername
	SecureServer.BindPRUDPEndPoint(SecureEndpoint)

	SecureServer.LibraryVersions.SetDefault(nex.NewLibraryVersion(globals.NEXMajor, globals.NEXMinor, globals.NEXPatch))
	SecureServer.AccessKey = globals.AccessKey
	// See authentication.go.
	SecureServer.ByteStreamSettings.UseStructureHeader = true

	SecureEndpoint.OnData(func(packet nex.PacketInterface) {
		request := packet.RMCMessage()
		if request == nil {
			return
		}
		pid := uint64(packet.Sender().PID())
		fmt.Printf("[MN9 Secure] PID=%d protocol=0x%02X method=0x%02X\n", pid, request.ProtocolID, request.MethodID)
	})

	SecureEndpoint.OnConnectionEnded(func(connection *nex.PRUDPConnection) {
		fmt.Printf("[MN9 Secure] PID=%d disconnected\n", uint64(connection.PID()))
	})

	globals.MatchmakingManager = common_globals.NewMatchmakingManager(SecureEndpoint, database.Postgres)
	globals.MatchmakingManager.GetUserFriendPIDs = globals.GetUserFriendPIDs

	registerSecureServerProtocols()

	port, _ := strconv.Atoi(os.Getenv("PN_MN9_SECURE_PORT"))
	globals.Logger.Successf("[MN9] Secure server listening on UDP %d", port)
	SecureServer.Listen(port)
}

// registerSecureServerProtocols: handshake, utility, Ranking, DataStore and
// the matchmaking stack (see the end of the function).
func registerSecureServerProtocols() {
	secureProtocol := secure.NewProtocol()
	SecureEndpoint.RegisterServiceProtocol(secureProtocol)
	secureCommon := common_secure.NewCommonProtocol(secureProtocol)
	secureCommon.EnableInsecureRegister()
	secureCommon.CreateReportDBRecord = database.CreateReportDBRecord

	utilityProtocol := utility.NewProtocol()
	SecureEndpoint.RegisterServiceProtocol(utilityProtocol)
	common_utility.NewCommonProtocol(utilityProtocol)

	rankingProtocol := ranking.NewProtocol()
	SecureEndpoint.RegisterServiceProtocol(rankingProtocol)
	rankingCommon := common_ranking.NewCommonProtocol(rankingProtocol)
	rankingCommon.GetRankingsAndCountByCategoryAndRankingOrderParam = database.MN9GetRankingsAndCountByCategoryAndRankingOrderParam
	rankingCommon.GetRankingsByMode = database.MN9GetRankings
	rankingCommon.GetCommonData = database.MN9GetCommonData
	rankingCommon.UploadCommonData = database.MN9UploadCommonData
	rankingCommon.InsertRankingByPIDAndRankingScoreData = database.MN9InsertRankingByPIDAndRankingScoreData

	datastoreProtocol := datastore.NewProtocol()
	SecureEndpoint.RegisterServiceProtocol(datastoreProtocol)
	datastoreCommon := common_datastore.NewCommonProtocol(datastoreProtocol)
	datastoreCommon.GetObjectInfosByDataStoreSearchParam = database.GetObjectInfosByDataStoreSearchParam
	datastoreCommon.InitializeObjectByPreparePostParam = database.InitializeObjectByPreparePostParam
	datastoreCommon.InitializeObjectRatingWithSlot = database.InitializeObjectRatingWithSlot
	datastoreCommon.GetObjectInfoByDataID = database.GetObjectInfoByDataID
	datastoreCommon.UpdateObjectPeriodByDataIDWithPassword = database.UpdateObjectPeriodByDataIDWithPassword
	datastoreCommon.UpdateObjectMetaBinaryByDataIDWithPassword = database.UpdateObjectMetaBinaryByDataIDWithPassword
	datastoreCommon.UpdateObjectDataTypeByDataIDWithPassword = database.UpdateObjectDataTypeByDataIDWithPassword
	datastoreCommon.GetObjectInfoByDataIDWithPassword = database.GetObjectInfoByDataIDWithPassword
	datastoreCommon.GetObjectInfoByPersistenceTargetWithPassword = database.GetObjectInfoByPersistenceTargetWithPassword
	datastoreProtocol.GetRatings = database.MN9GetRatings
	// Required by DataStore::CompletePostObject (last step of the score
	// upload/attachment flow).
	datastoreCommon.GetObjectOwnerByDataID = database.GetObjectOwnerByDataID
	datastoreCommon.GetObjectSizeByDataID = database.GetObjectSizeByDataID
	datastoreCommon.UpdateObjectUploadCompletedByDataID = database.UpdateObjectUploadCompletedByDataID
	datastoreCommon.DeleteObjectByDataID = database.DeleteObjectByDataID

	// DataStore::PreparePostObject (score-upload attachment) needs an
	// S3-compatible presigned-URL backend. Point PN_S3_ENDPOINT at any
	// S3-compatible service (self-hosted MinIO works well).
	s3Endpoint := os.Getenv("PN_S3_ENDPOINT")
	if s3Endpoint != "" {
		minioClient, err := minio.New(s3Endpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(os.Getenv("PN_S3_ACCESS_KEY"), os.Getenv("PN_S3_SECRET_KEY"), ""),
			Secure: true,
		})
		if err != nil {
			globals.Logger.Errorf("[MN9] Failed to create MinIO client: %s", err.Error())
		} else {
			s3Bucket := os.Getenv("PN_S3_BUCKET")
			if s3Bucket == "" {
				s3Bucket = "mn9-datastore"
			}
			datastoreCommon.S3Bucket = s3Bucket
			datastoreCommon.SetDataKeyBase("mn9")
			datastoreCommon.SetMinIOClient(minioClient)
		}
	} else {
		globals.Logger.Warning("[MN9] PN_S3_ENDPOINT not set - DataStore::PreparePostObject will fail")
	}

	// Step 2: the game calls MatchmakeExtension::GetFriendNotificationData
	// right after Register (seen on hardware), so the matchmaking stack is
	// registered with the stock common handlers.
	natTraversalProtocol := nat_traversal.NewProtocol()
	SecureEndpoint.RegisterServiceProtocol(natTraversalProtocol)
	common_nat_traversal.NewCommonProtocol(natTraversalProtocol)

	matchMakingProtocol := matchmaking.NewProtocol()
	SecureEndpoint.RegisterServiceProtocol(matchMakingProtocol)
	commonMatchMakingProtocol := common_matchmaking.NewCommonProtocol(matchMakingProtocol)
	commonMatchMakingProtocol.SetManager(globals.MatchmakingManager)

	matchMakingExtProtocol := matchmaking_ext.NewProtocol()
	SecureEndpoint.RegisterServiceProtocol(matchMakingExtProtocol)
	commonMatchMakingExtProtocol := common_matchmaking_ext.NewCommonProtocol(matchMakingExtProtocol)
	commonMatchMakingExtProtocol.SetManager(globals.MatchmakingManager)

	matchmakeExtensionProtocol := matchmake_extension.NewProtocol()
	SecureEndpoint.RegisterServiceProtocol(matchmakeExtensionProtocol)
	commonMatchmakeExtensionProtocol := common_matchmake_extension.NewCommonProtocol(matchmakeExtensionProtocol)
	commonMatchmakeExtensionProtocol.SetManager(globals.MatchmakingManager)
	// AutoMatchmake*Postpone hard-fail with Core::NotImplemented if these are nil.
	commonMatchmakeExtensionProtocol.CleanupMatchmakeSessionSearchCriterias = func(searchCriterias types.List[matchmaking_types.MatchmakeSessionSearchCriteria]) {}
	commonMatchmakeExtensionProtocol.CleanupSearchMatchmakeSession = func(matchmakeSession *matchmaking_types.MatchmakeSession) {}
}
