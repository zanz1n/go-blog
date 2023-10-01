const path = require("path");
const fs = require("fs");
const HtmlWebpackPlugin = require("html-webpack-plugin");
const CopyWebpackPlugin = require("copy-webpack-plugin");
const MiniCssExtractPlugin = require("mini-css-extract-plugin");
const CssMinimizerPlugin = require("css-minimizer-webpack-plugin");

const favicon = "favicon.svg";
const plugins = [];
/** @type {CopyWebpackPlugin.Pattern[]} */
const copies = [];

/**
 * @param {RegExp} regex 
 */
function loadTemplates(regex) {
    const entries = fs.readdirSync(path.join(__dirname, "templates"));

    for (let entry of entries) {
        if (regex.test(entry)) {
            const filenameS = entry.split(".");
            filenameS.pop();
            const filename = filenameS.join(".");

            plugins.push(new HtmlWebpackPlugin({
                template: path.join(__dirname, "templates", entry),
                filename: path.join("templates", entry),
                minify: false,
                chunks: [filename],
                publicPath: "/assets",
                favicon,
            }));
        } else {
            copies.push({
                from: path.join(__dirname, "templates", entry),
                to: path.join("templates", entry)
            });
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
