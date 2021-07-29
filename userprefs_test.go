package main

import (
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	prefsContent = []byte(`
// Mozilla User Preferences

// DO NOT EDIT THIS FILE.
//
// If you make changes to this file while the application is running,
// the changes will be overwritten when the application exits.
//
// To change a preference value, you can either:
// - modify it via the UI (e.g. via about:config in the browser); or
// - set it within a user.js file in your profile.

user_pref("app.normandy.first_run", false);
user_pref("app.normandy.migrationsApplied", 12);
user_pref("app.normandy.startupRolloutPrefs.browser.partnerlink.useAttributionURL", true);
user_pref("app.normandy.startupRolloutPrefs.browser.topsites.contile.enabled", true);
user_pref("app.normandy.startupRolloutPrefs.browser.topsites.experiment.ebay-2020-1", true);
user_pref("app.normandy.startupRolloutPrefs.browser.topsites.useRemoteSetting", true);
user_pref("app.normandy.user_id", "ca21ac2d-9836-4135-8643-90baefe2bcce");
user_pref("app.update.lastUpdateTime.addon-background-update-timer", 1626868179);
user_pref("app.update.lastUpdateTime.browser-cleanup-thumbnails", 1626957215);
user_pref("app.update.lastUpdateTime.recipe-client-addon-run", 1626867939);
user_pref("app.update.lastUpdateTime.region-update-timer", 1626627835);
user_pref("app.update.lastUpdateTime.rs-experiment-loader-timer", 1626957492);
user_pref("app.update.lastUpdateTime.search-engine-update-timer", 1626957335);
user_pref("app.update.lastUpdateTime.services-settings-poll-changes", 1626868059);
user_pref("app.update.lastUpdateTime.telemetry_modules_ping", 1626624835);
user_pref("app.update.lastUpdateTime.xpi-signature-verification", 1626868299);
user_pref("browser.aboutConfig.showWarning", false);
user_pref("browser.bookmarks.addedImportButton", true);
user_pref("browser.bookmarks.restore_default_bookmarks", false);
user_pref("browser.compactmode.show", true);
user_pref("browser.contentblocking.category", "standard");
user_pref("browser.contextual-services.contextId", "{9974547d-fdc8-4671-866c-b1138fcfca4b}");
user_pref("browser.download.panel.shown", true);
user_pref("browser.download.viewableInternally.typeWasRegistered.svg", true);
user_pref("browser.download.viewableInternally.typeWasRegistered.webp", true);
user_pref("browser.download.viewableInternally.typeWasRegistered.xml", true);
user_pref("browser.laterrun.bookkeeping.profileCreationTime", 1625838587);
user_pref("browser.laterrun.bookkeeping.sessionCount", 28);
user_pref("browser.laterrun.enabled", true);
user_pref("browser.migration.version", 109);
user_pref("browser.newtabpage.activity-stream.impressionId", "{3238e08e-13c5-4a46-8248-d6e73a0a7cea}");
user_pref("browser.newtabpage.pinned", "[]");
user_pref("browser.newtabpage.storageVersion", 1);
user_pref("browser.pageActions.persistedActions", "{\"ids\":[\"bookmark\"],\"idsInUrlbar\":[\"bookmark\"],\"idsInUrlbarPreProton\":[],\"version\":1}");
user_pref("browser.pagethumbnails.storage_version", 3);
user_pref("browser.proton.toolbar.version", 3);
user_pref("browser.region.update.updated", 1626627836);
user_pref("browser.rights.3.shown", true);
user_pref("browser.safebrowsing.provider.mozilla.lastupdatetime", "1626957187007");
user_pref("browser.safebrowsing.provider.mozilla.nextupdatetime", "1626978787007");
user_pref("browser.search.region", "FR");
user_pref("browser.sessionstore.upgradeBackup.latestBuildID", "20210714020445");
user_pref("browser.shell.defaultBrowserCheckCount", 15);
user_pref("browser.shell.didSkipDefaultBrowserCheckOnFirstRun", true);
user_pref("browser.startup.homepage_override.buildID", "20210714020445");
user_pref("browser.startup.homepage_override.mstone", "90.0");
user_pref("browser.startup.lastColdStartupCheck", 1626957491);
user_pref("browser.startup.upgradeDialog.version", 89);
user_pref("browser.tabs.tabClipWidth", 83);
user_pref("browser.tabs.tabMinWith", 94);
user_pref("browser.tabs.warnOnClose", false);
user_pref("browser.theme.toolbar-theme", 0);
user_pref("browser.toolbars.bookmarks.visibility", "always");
user_pref("browser.uiCustomization.state", "{\"placements\":{\"widget-overflow-fixed-list\":[],\"nav-bar\":[\"back-button\",\"forward-button\",\"stop-reload-button\",\"customizableui-special-spring1\",\"urlbar-container\",\"customizableui-special-spring2\",\"save-to-pocket-button\",\"downloads-button\",\"fxa-toolbar-menu-button\"],\"toolbar-menubar\":[\"menubar-items\"],\"TabsToolbar\":[\"tabbrowser-tabs\",\"new-tab-button\",\"alltabs-button\"],\"PersonalToolbar\":[\"import-button\",\"personal-bookmarks\"]},\"seen\":[\"save-to-pocket-button\",\"developer-button\"],\"dirtyAreaCache\":[\"nav-bar\",\"PersonalToolbar\",\"toolbar-menubar\",\"TabsToolbar\"],\"currentVersion\":17,\"newElementCount\":2}");
user_pref("browser.urlbar.placeholderName", "DuckDuckGo");
user_pref("browser.urlbar.placeholderName.private", "DuckDuckGo");
user_pref("browser.urlbar.suggest.calculator", true);
user_pref("browser.urlbar.tipShownCount.searchTip_onboard", 4);
user_pref("datareporting.policy.dataSubmissionPolicyAcceptedVersion", 2);
user_pref("datareporting.policy.dataSubmissionPolicyNotifiedTime", "1625838590028");
user_pref("devtools.everOpened", true);
user_pref("devtools.toolbox.splitconsoleEnabled", true);
user_pref("devtools.toolsidebar-height.inspector", 350);
user_pref("devtools.toolsidebar-width.inspector", 700);
user_pref("devtools.toolsidebar-width.inspector.splitsidebar", 350);
user_pref("distribution.Manjaro.bookmarksProcessed", true);
user_pref("distribution.iniFile.exists.appversion", "90.0");
user_pref("distribution.iniFile.exists.value", true);
user_pref("doh-rollout.balrog-migration-done", true);
user_pref("doh-rollout.doneFirstRun", true);
user_pref("doh-rollout.home-region", "FR");
user_pref("dom.push.userAgentID", "61d9a104140946528e51aa88fd08b343");
user_pref("extensions.activeThemeID", "{e0de5ee2-4619-413a-8300-a43a90196a6d}");
user_pref("extensions.blocklist.pingCountVersion", -1);
user_pref("extensions.databaseSchema", 33);
user_pref("extensions.getAddons.cache.lastUpdate", 1626868179);
user_pref("extensions.getAddons.databaseSchema", 6);
user_pref("extensions.incognito.migrated", true);
user_pref("extensions.lastAppBuildId", "20210714020445");
user_pref("extensions.lastAppVersion", "90.0");
user_pref("extensions.lastPlatformVersion", "90.0");
user_pref("extensions.pendingOperations", false);
user_pref("extensions.pictureinpicture.enable_picture_in_picture_overrides", true);
user_pref("extensions.reset_default_search.runonce.3", true);
user_pref("extensions.reset_default_search.runonce.reason", "previousRun");
user_pref("extensions.systemAddonSet", "{\"schema\":1,\"directory\":\"{49d78b53-bd0f-419d-9312-6068b0f686a2}\",\"addons\":{\"reset-search-defaults@mozilla.com\":{\"version\":\"2.0.0\"}}}");
user_pref("extensions.webcompat.enable_shims", true);
user_pref("extensions.webcompat.perform_injections", true);
user_pref("extensions.webcompat.perform_ua_overrides", true);
user_pref("extensions.webextensions.ExtensionStorageIDB.migrated.screenshots@mozilla.org", true);
user_pref("extensions.webextensions.uuids", "{\"doh-rollout@mozilla.org\":\"b9ecb8e7-691f-4134-8f98-105695f56aec\",\"formautofill@mozilla.org\":\"74a7def2-5ad4-4e34-a476-3b51818005f4\",\"pictureinpicture@mozilla.org\":\"bc991a02-d099-47c9-be28-2ecc10b7f7c2\",\"screenshots@mozilla.org\":\"7884d50d-6e0d-4f62-9bd8-91a964cf834a\",\"webcompat-reporter@mozilla.org\":\"210e0628-0f9e-42ae-a7cd-8b977803e6b8\",\"webcompat@mozilla.org\":\"eadb24f2-d0ab-4c7a-9026-c3ff60e721fb\",\"default-theme@mozilla.org\":\"b22d3a56-b8fe-452f-958f-d221d785395e\",\"google@search.mozilla.org\":\"66f66497-c162-4aa9-b096-938d49597de1\",\"wikipedia@search.mozilla.org\":\"080cb968-31e0-43ef-a84b-31f2c316a8df\",\"bing@search.mozilla.org\":\"29839e0b-ccb7-4ab6-849c-21e227b395bc\",\"ddg@search.mozilla.org\":\"e59e65ab-75df-46be-9a34-07a362d0e34a\",\"amazon@search.mozilla.org\":\"44a9435f-b258-4416-9824-ac2499e15fc4\",\"reset-search-defaults@mozilla.com\":\"d96dedb5-14b3-4175-b2a3-29a942302063\",\"{e0de5ee2-4619-413a-8300-a43a90196a6d}\":\"8b02f53b-0289-4b1d-895e-6c9f80541548\"}");
user_pref("fission.experiment.max-origins.last-disqualified", 0);
user_pref("fission.experiment.max-origins.last-qualified", 1625838590);
user_pref("fission.experiment.max-origins.qualified", true);
user_pref("gfx.webrender.all", true);
user_pref("idle.lastDailyNotification", 1626871507);
user_pref("layers.acceleration.force-enabled", true);
user_pref("layout.css.backdrop-filter.enabled", true);
user_pref("layout.css.color-mix.enabled", true);
user_pref("materialFox.reduceTabOverflow", true);
user_pref("media.gmp-gmpopenh264.abi", "x86_64-gcc3");
user_pref("media.gmp-gmpopenh264.lastUpdate", 1625840207);
user_pref("media.gmp-gmpopenh264.version", "1.8.1.1");
user_pref("media.gmp-manager.buildID", "20210714020445");
user_pref("media.gmp-manager.lastCheck", 1626867650);
user_pref("media.gmp.storage.version.observed", 1);
user_pref("network.trr.blocklist_cleanup_done", true);
user_pref("nimbus.syncdefaultsstore.upgradeDialog", "{\"slug\":\"upgradeDialog-defaultEnabled\",\"enabled\":true,\"targeting\":\"true\",\"variables\":{},\"description\":\"Turn on upgradeDialog by default for all users\"}");
user_pref("pdfjs.enabledCache.state", false);
user_pref("pdfjs.migrationVersion", 2);
user_pref("places.database.lastMaintenance", 1626626035);
user_pref("privacy.purge_trackers.date_in_cookie_database", "0");
user_pref("privacy.purge_trackers.last_purge", "1626871507933");
user_pref("privacy.sanitize.pending", "[]");
user_pref("security.insecure_connection_text.enabled", true);
user_pref("security.remote_settings.crlite_filters.checked", 1626883472);
user_pref("security.remote_settings.intermediates.checked", 1626882572);
user_pref("security.sandbox.content.tempDirSuffix", "0e4d920d-ffba-4a07-a1ab-2aa6389b2943");
user_pref("services.blocklist.addons-mlbf.checked", 1626957189);
user_pref("services.blocklist.gfx.checked", 1626957189);
user_pref("services.settings.clock_skew_seconds", -2);
user_pref("services.settings.last_etag", "\"1626944256233\"");
user_pref("services.settings.last_update_seconds", 1626957189);
user_pref("services.settings.main.anti-tracking-url-decoration.last_check", 1626957189);
user_pref("services.settings.main.cfr-fxa.last_check", 1626957189);
user_pref("services.settings.main.cfr.last_check", 1626957189);
user_pref("services.settings.main.doh-config.last_check", 1626957189);
user_pref("services.settings.main.doh-providers.last_check", 1626957189);
user_pref("services.settings.main.fxmonitor-breaches.last_check", 1626957189);
user_pref("services.settings.main.hijack-blocklists.last_check", 1626957189);
user_pref("services.settings.main.language-dictionaries.last_check", 1626957189);
user_pref("services.settings.main.message-groups.last_check", 1626957189);
user_pref("services.settings.main.nimbus-desktop-defaults.last_check", 1626957189);
user_pref("services.settings.main.nimbus-desktop-experiments.last_check", 1626957189);
user_pref("services.settings.main.normandy-recipes-capabilities.last_check", 1626957189);
user_pref("services.settings.main.partitioning-exempt-urls.last_check", 1626957189);
user_pref("services.settings.main.password-recipes.last_check", 1626957189);
user_pref("services.settings.main.pioneer-study-addons-v1.last_check", 1626957189);
user_pref("services.settings.main.public-suffix-list.last_check", 1626957189);
user_pref("services.settings.main.search-config.last_check", 1626957189);
user_pref("services.settings.main.search-default-override-allowlist.last_check", 1626957189);
user_pref("services.settings.main.search-telemetry.last_check", 1626957189);
user_pref("services.settings.main.sites-classification.last_check", 1626957189);
user_pref("services.settings.main.tippytop.last_check", 1626957189);
user_pref("services.settings.main.top-sites.last_check", 1626957189);
user_pref("services.settings.main.url-classifier-skip-urls.last_check", 1626957189);
user_pref("services.settings.main.websites-with-shared-credential-backends.last_check", 1626957189);
user_pref("services.settings.main.whats-new-panel.last_check", 1626957189);
user_pref("services.settings.security.onecrl.checked", 1626882572);
user_pref("services.sync.clients.lastSync", "0");
user_pref("services.sync.declinedEngines", "");
user_pref("services.sync.globalScore", 0);
user_pref("services.sync.nextSync", 0);
user_pref("services.sync.tabs.lastSync", "0");
user_pref("storage.vacuum.last.index", 1);
user_pref("storage.vacuum.last.places.sqlite", 1626626035);
user_pref("svg.context-properties.content.enabled", true);
user_pref("toolkit.legacyUserProfileCustomizations.stylesheets", true);
user_pref("toolkit.startup.last_success", 1626957490);
user_pref("toolkit.telemetry.cachedClientID", "701ca6a7-c501-47d6-aef9-2a1d57a8d1eb");
user_pref("toolkit.telemetry.pioneer-new-studies-available", true);
user_pref("toolkit.telemetry.previousBuildID", "20210714020445");
user_pref("toolkit.telemetry.reportingpolicy.firstRun", false);
user_pref("trailhead.firstrun.didSeeAboutWelcome", true);
	`)
)

func TestToUserJS(t *testing.T) {
	value, err := ToUserJSFile(map[string]interface{}{
		"browser.tabs.tabClipWidth":              90,
		"svg.context-properties.content.enabled": true,
	})
	if err != nil {
		panic(err)
	}
	assert.Equal(t,
		`user_pref("browser.tabs.tabClipWidth", 90);
user_pref("svg.context-properties.content.enabled", true);`,
		sortLines(value),
	)
}

func sortLines(s string) string {
	lines := strings.Split(s, "\n")
	sort.Strings(lines)
	return strings.Join(lines, "\n")
}

func TestValueOfUserPrefCall(t *testing.T) {
	result, err := ValueOfUserPrefCall(prefsContent, "storage.vacuum.last.index")
	assert.NoError(t, err)
	assert.Equal(t, "1", result)

	result, err = ValueOfUserPrefCall(prefsContent, "app.normandy.first_run")
	assert.NoError(t, err)
	assert.Equal(t, "false", result)

	result, err = ValueOfUserPrefCall(prefsContent, "browser.search.region")
	assert.NoError(t, err)
	assert.Equal(t, "FR", result)

	result, err = ValueOfUserPrefCall(prefsContent, "lkghjoertkjhoietrjhoirtjhoirtjhor")
	assert.Contains(t, err.Error(), `key "lkghjoertkjhoietrjhoirtjhoirtjhor" not found`)
	assert.Equal(t, result, "")
}
