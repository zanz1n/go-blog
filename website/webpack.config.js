const path = require("path");
const fs = require("fs");
const HtmlWebpackPlugin = require("html-webpack-plugin");
const CopyWebpackPlugin = require("copy-webpack-plugin");
const MiniCssExtractPlugin = require("mini-css-extract-plugin");
const CssMinimizerPlugin = require("css-minimizer-webpack-plugin");

const favicon = "favicon.svg";
const minifyTemplates = true;

const plugins = [];
/** @type {CopyWebpackPlugin.Pattern[]} */
const copies = [];

/** @param {RegExp} regex */
function loadTemplates(regex, sub = false, dir = "templates") {
    const entries = fs.readdirSync(path.join(__dirname, dir));

    for (let entry of entries) {
        if (regex.test(entry)) {
            const filenameS = entry.split(".");
            filenameS.pop();
            const filename = filenameS.join(".");

            if (!sub) {
                plugins.push(new HtmlWebpackPlugin({
                    template: path.join(__dirname, dir, entry),
                    filename: path.join(dir, entry),
                    minify: minifyTemplates,
                    chunks: [filename],
                    publicPath: "/assets",
                    favicon
                }));
            } else {
                plugins.push(new HtmlWebpackPlugin({
                    template: path.join(__dirname, dir, entry),
                    filename: path.join(dir, entry),
                    minify: minifyTemplates,
                    chunks: [],
                }));
            }
        } else {
            const fstat = fs.statSync(path.join(__dirname, dir, entry));
            if (fstat.isDirectory()) {
                loadTemplates(regex, true, path.join(dir, entry));
            } else {
                copies.push({
                    from: path.join(__dirname, dir, entry),
                    to: path.join(dir, entry)
                });
            }
        }
    }
}

function loadScriptEntries(ext = ".js") {
    const entries = fs.readdirSync(path.join(__dirname, "src"));

    const res = {};

    for (let entry of entries) {
        if (entry.endsWith(ext) && entry.length > ext.length) {
            const stat = fs.statSync(path.join(__dirname, "src", entry));

            if (stat.isFile()) {
                const entryKey = entry.replace(ext, "");
                res[entryKey] = "./src/" + entry;
            }
        }
    }

    return res;
}

function loadPublicAssets(dir = "./public") {
    const entries = fs.readdirSync(path.join(__dirname, dir));

    for (let entry of entries) {
        copies.push({
            from: path.join(__dirname, dir, entry),
            to: entry
        });
    }
}

// CSS
plugins.push(new MiniCssExtractPlugin({
    filename: "[name]-[contenthash].css"
}));

// Templates
loadTemplates(/\.hbs$/);
if (copies.length > 0) {
    plugins.push(new CopyWebpackPlugin({
        patterns: copies,
    }));
}

// Assets
loadPublicAssets();

/**@type {import("webpack").Configuration} */
module.exports = {
    plugins,

    mode: "production",

    entry: loadScriptEntries(".ts"),

    module: {
        rules: [
            {
                test: /\.ts$/,
                exclude: /node_modules/,
                use: ["ts-loader"]
            },
            {
                test: /\.css$/,
                exclude: /node_modules/,
                use: [
                    MiniCssExtractPlugin.loader,
                    "css-loader"
                ]
            }
        ],
    },

    resolve: {
        extensions: [".ts", ".js", ".css"],
    },

    output: {
        path: path.join(__dirname, "dist"),
        filename: "[name]-[contenthash].js",
        clean: true
    },

    optimization: {
        minimize: true,
        minimizer: [
            "...",
            new CssMinimizerPlugin()
        ]
    }
};
