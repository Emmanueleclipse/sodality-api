package ipfs

// import (
// 	"context"
// 	"flag"
// 	"fmt"
// 	"io/ioutil"
// 	"log"
// 	"os"
// 	"path/filepath"
// 	"strings"
// 	"sync"

// 	config "github.com/ipfs/go-ipfs-config"
// 	files "github.com/ipfs/go-ipfs-files"

// 	"github.com/ipfs/go-ipfs/core"
// 	"github.com/ipfs/go-ipfs/core/coreapi"
// 	"github.com/ipfs/go-ipfs/core/node/libp2p"
// 	"github.com/ipfs/go-ipfs/plugin/loader"
// 	"github.com/ipfs/go-ipfs/repo/fsrepo"
// 	icore "github.com/ipfs/interface-go-ipfs-core"
// 	"github.com/ipfs/interface-go-ipfs-core/path"
// 	"github.com/libp2p/go-libp2p-core/peer"
// 	ma "github.com/multiformats/go-multiaddr"
// )

// var flagExp = flag.Bool("experimental", false, "enable experimental features")

// func StartingIPFS(ctx context.Context) icore.CoreAPI {
// 	log.Println("-- Getting an IPFS node running -- ")
// 	log.Println("Spawning node on a temporary repo")
// 	ipfs, err := SpawnEphemeral(ctx)
// 	if err != nil {
// 		log.Fatalln(fmt.Errorf("failed to spawn ephemeral node: %s", err))
// 	}
// 	return ipfs
// }

// func AddFilePath(file []byte, ctx context.Context, ipfs icore.CoreAPI) (path.Resolved, error) {

// 	cidFile, err := ipfs.Unixfs().Add(ctx, files.NewBytesFile(file))
// 	if err != nil {
// 		return nil, fmt.Errorf("unable to add file: %v", err)
// 	}

// 	return cidFile, nil
// }

// func GetFile(cidFile path.Path, ctx context.Context, ipfs icore.CoreAPI) (files.Node, string, error) {
// 	outputBasePath, err := ioutil.TempDir("", "example")
// 	if err != nil {
// 		return nil, "", fmt.Errorf("unable to create output dir: %v", err)
// 	}

// 	outputPathFile := outputBasePath + "/" + strings.Split(cidFile.String(), "/")[2]
// 	rootNodeFile, err := ipfs.Unixfs().Get(ctx, cidFile)
// 	if err != nil {
// 		return nil, "", fmt.Errorf("unable to getfile with cid: %v", err)
// 	}

// 	err = files.WriteTo(rootNodeFile, outputPathFile)
// 	if err != nil {
// 		return nil, "", fmt.Errorf("unable to fetched cid: %v", err)
// 	}

// 	return rootNodeFile, outputPathFile, nil
// }

// func GetUnixfsNode(path string) (files.Node, error) {
// 	st, err := os.Stat(path)
// 	if err != nil {
// 		return nil, err
// 	}

// 	f, err := files.NewSerialFile(path, false, st)

// 	if err != nil {
// 		return nil, err
// 	}

// 	return f, nil
// }

// var loadPluginsOnce sync.Once

// func SpawnEphemeral(ctx context.Context) (icore.CoreAPI, error) {
// 	var onceErr error
// 	loadPluginsOnce.Do(func() {
// 		onceErr = SetupPlugins("")
// 	})
// 	if onceErr != nil {
// 		return nil, onceErr
// 	}
// 	// Create a Temporary Repo
// 	repoPath, err := CreateTempRepo()
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to create temp repo: %s", err)
// 	}

// 	// Spawning an ephemeral IPFS node
// 	return CreateNode(ctx, repoPath)
// }

// var nodee *core.IpfsNode

// func CreateNode(ctx context.Context, repoPath string) (icore.CoreAPI, error) {
// 	// Open the repo
// 	repo, err := fsrepo.Open(repoPath)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Construct the node

// 	nodeOptions := &core.BuildCfg{
// 		Online:  true,
// 		Routing: libp2p.DHTOption, // This option sets the node to be a full DHT node (both fetching and storing DHT Records)
// 		// Routing: libp2p.DHTClientOption, // This option sets the node to be a client DHT node (only fetching records)
// 		Repo: repo,
// 	}

// 	nodee, err = core.NewNode(ctx, nodeOptions)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Attach the Core API to the constructed node
// 	return coreapi.NewCoreAPI(nodee)
// }

// func CreateTempRepo() (string, error) {
// 	repoPath, err := ioutil.TempDir("", "ipfs-shell")
// 	if err != nil {
// 		return "", fmt.Errorf("failed to get temp dir: %s", err)
// 	}

// 	// Create a config with default options and a 2048 bit key
// 	cfg, err := config.Init(ioutil.Discard, 2048)
// 	if err != nil {
// 		return "", err
// 	}

// 	// When creating the repository, you can define custom settings on the repository, such as enabling experimental
// 	// features (See experimental-features.md) or customizing the gateway endpoint.
// 	// To do such things, you should modify the variable `cfg`. For example:
// 	if *flagExp {
// 		// https://github.com/ipfs/go-ipfs/blob/master/docs/experimental-features.md#ipfs-filestore
// 		cfg.Experimental.FilestoreEnabled = true
// 		// https://github.com/ipfs/go-ipfs/blob/master/docs/experimental-features.md#ipfs-urlstore
// 		cfg.Experimental.UrlstoreEnabled = true
// 		// https://github.com/ipfs/go-ipfs/blob/master/docs/experimental-features.md#directory-sharding--hamt
// 		// cfg.Experimental.ShardingEnabled = true
// 		// https://github.com/ipfs/go-ipfs/blob/master/docs/experimental-features.md#ipfs-p2p
// 		cfg.Experimental.Libp2pStreamMounting = true
// 		// https://github.com/ipfs/go-ipfs/blob/master/docs/experimental-features.md#p2p-http-proxy
// 		cfg.Experimental.P2pHttpProxy = true
// 		// https://github.com/ipfs/go-ipfs/blob/master/docs/experimental-features.md#strategic-providing
// 		cfg.Experimental.StrategicProviding = true
// 	}

// 	// repoPath := ".ipfs"
// 	// _, err = fsrepo.Open(repoPath)
// 	// if err != nil {
// 	// 	err = fsrepo.Init(repoPath, cfg)
// 	// 	if err != nil {
// 	// 		return "", fmt.Errorf("failed to init ephemeral node: %s", err)
// 	// 	}
// 	// }

// 	// Create the repo with the config
// 	err = fsrepo.Init(repoPath, cfg)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to init ephemeral node: %s", err)
// 	}
// 	return repoPath, nil
// }

// func SetupPlugins(externalPluginsPath string) error {
// 	// Load any external plugins if available on externalPluginsPath
// 	plugins, err := loader.NewPluginLoader(filepath.Join(externalPluginsPath, "plugins"))
// 	if err != nil {
// 		return fmt.Errorf("error loading plugins: %s", err)
// 	}

// 	// Load preloaded and external plugins
// 	if err := plugins.Initialize(); err != nil {
// 		return fmt.Errorf("error initializing plugins: %s", err)
// 	}

// 	if err := plugins.Inject(); err != nil {
// 		return fmt.Errorf("error initializing plugins: %s", err)
// 	}

// 	return nil
// }

// func ConnectPeers(ctx context.Context, ipfs icore.CoreAPI) {
// 	// peerMa := fmt.Sprintf("/ip4/127.0.0.1/udp/4010/p2p/%s", nodee.Identity.String())

// 	bootstrapNodes := []string{
// 		// IPFS Bootstrapper nodes.

// 		"/dnsaddr/bootstrap.libp2p.io/p2p/QmNnooDu7bfjPFoTZYxMNLWUQJyrVwtbZg5gBMjTezGAJN",
// 		"/dnsaddr/bootstrap.libp2p.io/p2p/QmQCU2EcMqAqQPR2i9bChDtGNJchTbq5TbXJJ16u19uLTa",
// 		"/dnsaddr/bootstrap.libp2p.io/p2p/QmbLHAnMoJPWSCR5Zhtx6BHJX9KiKNN6tpvbUcqanj75Nb",
// 		"/dnsaddr/bootstrap.libp2p.io/p2p/QmcZf59bWwK5XFi76CZX8cbJ4BhTzzA3gU1ZjYZcYW3dwt",

// 		// IPFS Cluster Pinning nodes
// 		// "/ip4/13.250.48.222/tcp/4001/p2p/QmajoVpyBjGj84j3Pmvq2Svbo9Ec6maEKfpeXAhTmmTciY",
// 		"/ip4/138.201.67.219/tcp/4001/p2p/QmUd6zHcbkbcs7SMxwLs48qZVX3vpcM8errYS7xEczwRMA",
// 		"/ip4/138.201.67.219/udp/4001/quic/p2p/QmUd6zHcbkbcs7SMxwLs48qZVX3vpcM8errYS7xEczwRMA",
// 		"/ip4/138.201.67.220/tcp/4001/p2p/QmNSYxZAiJHeLdkBg38roksAR9So7Y5eojks1yjEcUtZ7i",
// 		"/ip4/138.201.67.220/udp/4001/quic/p2p/QmNSYxZAiJHeLdkBg38roksAR9So7Y5eojks1yjEcUtZ7i",
// 		"/ip4/138.201.68.74/tcp/4001/p2p/QmdnXwLrC8p1ueiq2Qya8joNvk3TVVDAut7PrikmZwubtR",
// 		"/ip4/138.201.68.74/udp/4001/quic/p2p/QmdnXwLrC8p1ueiq2Qya8joNvk3TVVDAut7PrikmZwubtR",
// 		"/ip4/94.130.135.167/tcp/4001/p2p/QmUEMvxS2e7iDrereVYc5SWPauXPyNwxcy9BXZrC1QTcHE",
// 		"/ip4/94.130.135.167/udp/4001/quic/p2p/QmUEMvxS2e7iDrereVYc5SWPauXPyNwxcy9BXZrC1QTcHE",

// 		// You can add more nodes here, for example, another IPFS node you might have running locally, mine was:
// 		"/ip4/127.0.0.1/tcp/4010/p2p/QmZp2fhDLxjYue2RiUvLwT9MWdnbDxam32qYFnGmxZDh5L",
// 		"/ip4/127.0.0.1/udp/4010/quic/p2p/QmZp2fhDLxjYue2RiUvLwT9MWdnbDxam32qYFnGmxZDh5L",
// 		// peerMa,
// 	}

// 	go func() {
// 		ConnectToPeers(ctx, ipfs, bootstrapNodes)
// 		// if err != nil {
// 		// 	log.Printf("failed connect to peers: %s", err)
// 		// }
// 	}()
// }

// func ConnectToPeers(ctx context.Context, ipfs icore.CoreAPI, peers []string) error {
// 	var wg sync.WaitGroup
// 	peerInfos := make(map[peer.ID]*peer.AddrInfo, len(peers))
// 	for _, addrStr := range peers {
// 		addr, err := ma.NewMultiaddr(addrStr)
// 		if err != nil {
// 			return err
// 		}
// 		pii, err := peer.AddrInfoFromP2pAddr(addr)
// 		if err != nil {
// 			return err
// 		}
// 		pi, ok := peerInfos[pii.ID]
// 		if !ok {
// 			pi = &peer.AddrInfo{ID: pii.ID}
// 			peerInfos[pi.ID] = pi
// 		}
// 		pi.Addrs = append(pi.Addrs, pii.Addrs...)
// 	}

// 	wg.Add(len(peerInfos))
// 	for _, peerInfo := range peerInfos {
// 		go func(peerInfo *peer.AddrInfo) {
// 			defer wg.Done()
// 			ipfs.Swarm().Connect(ctx, *peerInfo)
// 			// if err != nil {
// 			// 	log.Printf("failed to connect to %s: %s", peerInfo.ID, err)
// 			// }
// 		}(peerInfo)
// 	}
// 	wg.Wait()
// 	return nil
// }
