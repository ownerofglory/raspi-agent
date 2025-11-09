export default function Dashboard() {
    const devices = [
        {
            name: "Living Room Pi",
            model: "Raspberry Pi 4",
            status: "Online",
            color: "#28A745",
            info: "IP: 192.168.1.101",
        },
        {
            name: "Kitchen Assistant",
            model: "Raspberry Pi 3B+",
            status: "Offline",
            color: "#6C757D",
            info: "Last seen: 2 hours ago",
        },
        {
            name: "Workshop Pi",
            model: "Raspberry Pi 4",
            status: "Needs Attention",
            color: "#FFC107",
            info: "Storage: 95% full",
        },
        {
            name: "Garage Speaker",
            model: "Raspberry Pi Zero W",
            status: "Connecting",
            color: "#137fec",
            info: "Attempting to connect...",
            animate: true,
        },
    ];

    return (
        <div className="font-display bg-background-light dark:bg-background-dark text-text-light dark:text-text-dark min-h-screen flex flex-col">
            <div className="flex flex-1 justify-center px-4 py-5 sm:px-6 md:px-8 lg:px-10">
                <div className="w-full max-w-6xl flex flex-col">
                    {/* Top NavBar */}
                    <header className="flex items-center justify-between whitespace-nowrap border-b border-border-light dark:border-border-dark px-6 py-4">
                        <div className="flex items-center gap-4">
                            <div className="h-8 w-8 text-primary">
                                <svg
                                    fill="none"
                                    viewBox="0 0 48 48"
                                    xmlns="http://www.w3.org/2000/svg"
                                >
                                    <g clipPath="url(#clip0_6_319)">
                                        <path
                                            d="M8.57829 8.57829C5.52816 11.6284 3.451 15.5145 2.60947 19.7452C1.76794 23.9758 2.19984 28.361 3.85056 32.3462C5.50128 36.3314 8.29667 39.7376 11.8832 42.134C15.4698 44.5305 19.6865 45.8096 24 45.8096C28.3135 45.8096 32.5302 44.5305 36.1168 42.134C39.7033 39.7375 42.4987 36.3314 44.1494 32.3462C45.8002 28.361 46.2321 23.9758 45.3905 19.7452C44.549 15.5145 42.4718 11.6284 39.4217 8.57829L24 24L8.57829 8.57829Z"
                                            fill="currentColor"
                                        ></path>
                                    </g>
                                    <defs>
                                        <clipPath id="clip0_6_319">
                                            <rect fill="white" width="48" height="48" />
                                        </clipPath>
                                    </defs>
                                </svg>
                            </div>
                            <h2 className="text-xl font-bold">DIY Voice Assistant</h2>
                        </div>

                        <div className="flex items-center gap-4">
                            <IconButton icon="settings" />
                            <IconButton icon="notifications" />
                            <div
                                className="size-10 rounded-full bg-cover bg-center bg-no-repeat"
                                style={{
                                    backgroundImage:
                                        "url('https://lh3.googleusercontent.com/aida-public/AB6AXuDYfzo5koOwtooScNdf6fBAamCWgNF0BMCG7KEIDCU1dwHESe51RN752KIA85lu6itU-6z4pNCBOSd36sR8Kqu3ZOfHxRytlVv6NzrdtTMV4162Gf_BxxfAZYa1aBxWqRkoTx-MO5tg7o00KMzPei0FikdKFE6hlSbtjj6eI1WtmQYOhasQvUyY3yw9SaD1s3FXRIZR6aVDg4C1uoEangIYdQZdmQOcU-pEOn-rzFGDNlBuBdvscBKuKrb9O1lKV9GN9Ogxz747JtnR')",
                                }}
                            />
                        </div>
                    </header>

                    {/* Main Section */}
                    <main className="flex flex-col gap-8 p-4 pt-8 md:p-6">
                        {/* Heading */}
                        <div className="flex flex-wrap items-center justify-between gap-4">
                            <div className="flex min-w-72 flex-col gap-1">
                                <p className="text-4xl font-bold tracking-tighter">
                                    My Devices
                                </p>
                                <p className="text-base text-text-muted-light dark:text-text-muted-dark">
                                    Manage your registered Raspberry Pi voice assistants
                                </p>
                            </div>
                            <button className="flex h-12 min-w-[84px] max-w-[480px] cursor-pointer items-center justify-center gap-2 overflow-hidden rounded-lg bg-primary px-6 text-base font-bold text-white shadow-lg shadow-primary/30 transition-all hover:bg-primary/90">
                                <span className="material-symbols-outlined">add_circle</span>
                                <span className="truncate">Enroll a New Device</span>
                            </button>
                        </div>

                        {/* Device Grid */}
                        <div className="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
                            {devices.map((device, idx) => (
                                <DeviceCard key={idx} {...device} />
                            ))}
                        </div>
                    </main>
                </div>
            </div>
        </div>
    );
}

/* ========== Reusable Components ========== */

function IconButton({ icon }: { icon: string }) {
    return (
        <button className="flex h-10 w-10 cursor-pointer items-center justify-center overflow-hidden rounded-lg bg-transparent text-text-muted-light dark:text-text-muted-dark hover:bg-black/5 dark:hover:bg-white/5">
            <span className="material-symbols-outlined">{icon}</span>
        </button>
    );
}

interface DeviceCardProps {
    name: string;
    model: string;
    status: string;
    color: string;
    info: string;
    animate?: boolean;
}

function DeviceCard({
                        name,
                        model,
                        status,
                        color,
                        info,
                        animate,
                    }: DeviceCardProps) {
    return (
        <div className="flex flex-col gap-4 rounded-xl border border-border-light dark:border-border-dark bg-card-light dark:bg-card-dark p-5 shadow-sm transition-all hover:shadow-lg hover:-translate-y-1">
            <div className="flex items-center justify-between">
                <p className="text-lg font-bold">{name}</p>
                <div className="flex items-center gap-2">
                    <div
                        className={`h-2.5 w-2.5 rounded-full ${
                            animate ? "animate-pulse" : ""
                        }`}
                        style={{ backgroundColor: color }}
                    ></div>
                    <p className="text-sm font-medium" style={{ color }}>
                        {status}
                    </p>
                </div>
            </div>
            <div className="flex flex-col text-sm text-text-muted-light dark:text-text-muted-dark">
                <span>{model}</span>
                <span>{info}</span>
            </div>
            <button className="mt-2 flex h-10 w-full cursor-pointer items-center justify-center overflow-hidden rounded-lg bg-black/5 text-sm font-bold text-text-light dark:bg-white/10 dark:text-text-dark hover:bg-black/10 dark:hover:bg-white/20">
                <span className="truncate">Manage Device</span>
            </button>
        </div>
    );
}
