using System;
using System.IO;
using System.Text;
using System.Threading.Tasks;
using Windows.Data.Pdf;
using Windows.Graphics.Imaging;
using Windows.Media.Ocr;
using Windows.Storage;
using Windows.Storage.Streams;

namespace PdfTextExtractor
{
    class Program
    {
        static async Task<int> Main(string[] args)
        {
            try
            {
                // Parser les arguments
                if (args.Length < 1)
                {
                    Console.Error.WriteLine("Usage: PdfTextExtractor <pdf-file> [language] [dpi]");
                    Console.Error.WriteLine("Example: PdfTextExtractor document.pdf fr-FR 300");
                    return 1;
                }

                string pdfPath = args[0];
                string language = args.Length > 1 ? args[1] : "fr-FR";
                int dpi = args.Length > 2 ? int.Parse(args[2]) : 300;

                // Vérifier que le fichier existe
                if (!File.Exists(pdfPath))
                {
                    Console.Error.WriteLine($"Erreur: Le fichier n'existe pas: {pdfPath}");
                    return 1;
                }

                // Extraire le texte
                string text = await ExtractTextFromPdfAsync(pdfPath, language, dpi);

                // Afficher le résultat sur stdout
                Console.WriteLine("========== TEXTE EXTRAIT ==========");
                Console.WriteLine(text);
                Console.WriteLine("===================================");

                return 0;
            }
            catch (Exception ex)
            {
                Console.Error.WriteLine($"Erreur: {ex.Message}");
                Console.Error.WriteLine(ex.StackTrace);
                return 1;
            }
        }

        static async Task<string> ExtractTextFromPdfAsync(string pdfPath, string languageCode, int dpi)
        {
            var result = new StringBuilder();

            // Ouvrir le fichier PDF
            var file = await StorageFile.GetFileFromPathAsync(Path.GetFullPath(pdfPath));
            var pdfDocument = await PdfDocument.LoadFromFileAsync(file);

            Console.Error.WriteLine($"Traitement de {pdfDocument.PageCount} page(s)...");

            // Initialiser l'OCR engine
            var language = new Windows.Globalization.Language(languageCode);
            var ocrEngine = OcrEngine.TryCreateFromLanguage(language);

            if (ocrEngine == null)
            {
                // Essayer avec l'anglais par défaut
                Console.Error.WriteLine($"Langue {languageCode} non disponible, utilisation de l'anglais");
                language = new Windows.Globalization.Language("en-US");
                ocrEngine = OcrEngine.TryCreateFromLanguage(language);
            }

            if (ocrEngine == null)
            {
                throw new Exception("Impossible d'initialiser l'engine OCR");
            }

            Console.Error.WriteLine($"OCR initialisé: {ocrEngine.RecognizerLanguage.DisplayName}");

            // Créer un dossier temporaire
            string tempFolder = Path.Combine(Path.GetTempPath(), $"pdf2txt_{Guid.NewGuid():N}");
            Directory.CreateDirectory(tempFolder);

            try
            {
                // Traiter chaque page
                for (uint pageIndex = 0; pageIndex < pdfDocument.PageCount; pageIndex++)
                {
                    uint pageNumber = pageIndex + 1;
                    Console.Error.WriteLine($"Traitement de la page {pageNumber}/{pdfDocument.PageCount}...");

                    using (var pdfPage = pdfDocument.GetPage(pageIndex))
                    {
                        // Créer un fichier temporaire pour l'image
                        string imagePath = Path.Combine(tempFolder, $"page_{pageNumber:D4}.png");
                        var imageFile = await StorageFile.GetFileFromPathAsync(
                            Path.GetFullPath(imagePath)
                        ).AsTask().ContinueWith(async t =>
                        {
                            if (t.IsFaulted)
                            {
                                // Le fichier n'existe pas, le créer
                                var folder = await StorageFolder.GetFolderFromPathAsync(tempFolder);
                                return await folder.CreateFileAsync($"page_{pageNumber:D4}.png", CreationCollisionOption.ReplaceExisting);
                            }
                            return await t;
                        }).Unwrap();

                        // Rendre la page en image
                        var renderOptions = new PdfPageRenderOptions
                        {
                            DestinationWidth = (uint)(pdfPage.Size.Width * dpi / 72),
                            DestinationHeight = (uint)(pdfPage.Size.Height * dpi / 72)
                        };

                        using (var stream = await imageFile.OpenAsync(FileAccessMode.ReadWrite))
                        {
                            await pdfPage.RenderToStreamAsync(stream, renderOptions);
                            stream.Dispose();
                        }

                        // Faire l'OCR sur l'image
                        using (var stream = await imageFile.OpenAsync(FileAccessMode.Read))
                        {
                            var decoder = await BitmapDecoder.CreateAsync(stream);
                            var softwareBitmap = await decoder.GetSoftwareBitmapAsync();
                            var ocrResult = await ocrEngine.RecognizeAsync(softwareBitmap);

                            string pageText = ocrResult.Text;

                            if (!string.IsNullOrWhiteSpace(pageText))
                            {
                                result.AppendLine($"=== Page {pageNumber} ===");
                                result.AppendLine(pageText);
                                result.AppendLine();
                                Console.Error.WriteLine($"  Texte extrait: {pageText.Length} caractères");
                            }
                            else
                            {
                                Console.Error.WriteLine($"  Aucun texte trouvé");
                            }
                        }
                    }
                }
            }
            finally
            {
                // Nettoyer le dossier temporaire
                try
                {
                    Directory.Delete(tempFolder, true);
                }
                catch
                {
                    // Ignorer les erreurs de nettoyage
                }
            }

            return result.ToString();
        }
    }
}
