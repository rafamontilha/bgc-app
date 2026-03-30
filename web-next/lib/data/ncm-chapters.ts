import type { NcmChapter } from '@/lib/types/onboarding';

/**
 * Lista completa dos 97 capítulos NCM (Nomenclatura Comum do Mercosul).
 *
 * - `code`: 2 dígitos, com zero à esquerda.
 * - `label`: descrição oficial em português (simplificada para legibilidade no autocomplete).
 * - `defaultNcm8d`: NCM de 8 dígitos mais representativo do capítulo por volume
 *   de exportação brasileira (Comex Stat 2024). Usado para pré-carregar o simulador.
 *
 * Capítulo 77 reservado (não utilizado na NCM/SH vigente).
 */
export const NCM_CHAPTERS: NcmChapter[] = [
  // Seção I — Animais Vivos e Produtos do Reino Animal
  { code: '01', label: 'Animais vivos', defaultNcm8d: '01029090' },
  { code: '02', label: 'Carnes e miudezas comestíveis', defaultNcm8d: '02013000' },
  { code: '03', label: 'Peixes, crustáceos e outros invertebrados aquáticos', defaultNcm8d: '03034200' },
  { code: '04', label: 'Leite, laticínios, ovos de aves e mel natural', defaultNcm8d: '04021010' },
  { code: '05', label: 'Outros produtos de origem animal', defaultNcm8d: '05040010' },

  // Seção II — Produtos do Reino Vegetal
  { code: '06', label: 'Plantas vivas e produtos de floricultura', defaultNcm8d: '06024090' },
  { code: '07', label: 'Legumes, hortaliças, raízes e tubérculos comestíveis', defaultNcm8d: '07019000' },
  { code: '08', label: 'Frutas; cascas de frutas cítricas ou de melões', defaultNcm8d: '08051000' },
  { code: '09', label: 'Café, chá, mate e especiarias', defaultNcm8d: '09010110' },
  { code: '10', label: 'Cereais', defaultNcm8d: '10011900' },
  { code: '11', label: 'Produtos de moagem; malte; amidos e féculas; glúten de trigo', defaultNcm8d: '11010010' },
  { code: '12', label: 'Sementes e frutos oleaginosos; grãos e sementes diversas', defaultNcm8d: '12010090' },
  { code: '13', label: 'Gomas, resinas e outros sucos e extratos vegetais', defaultNcm8d: '13019090' },
  { code: '14', label: 'Matérias para entrançar e outros produtos de origem vegetal', defaultNcm8d: '14049090' },

  // Seção III — Gorduras e Óleos
  { code: '15', label: 'Gorduras e óleos animais ou vegetais', defaultNcm8d: '15079019' },

  // Seção IV — Produtos das Indústrias Alimentares
  { code: '16', label: 'Preparações de carne, peixes ou crustáceos', defaultNcm8d: '16010090' },
  { code: '17', label: 'Açúcares e produtos de confeitaria', defaultNcm8d: '17011400' },
  { code: '18', label: 'Cacau e suas preparações', defaultNcm8d: '18010000' },
  { code: '19', label: 'Preparações à base de cereais, farinha, amido ou leite', defaultNcm8d: '19041010' },
  { code: '20', label: 'Preparações de produtos hortícolas, fruta ou partes de plantas', defaultNcm8d: '20091100' },
  { code: '21', label: 'Preparações alimentícias diversas', defaultNcm8d: '21069090' },
  { code: '22', label: 'Bebidas, líquidos alcoólicos e vinagres', defaultNcm8d: '22021000' },
  { code: '23', label: 'Resíduos e alimentos para animais', defaultNcm8d: '23040010' },
  { code: '24', label: 'Tabaco e seus sucedâneos', defaultNcm8d: '24011020' },

  // Seção V — Produtos Minerais
  { code: '25', label: 'Sal; enxofre; terras e pedras; gesso, cal e cimento', defaultNcm8d: '25010019' },
  { code: '26', label: 'Minérios, escórias e cinzas', defaultNcm8d: '26011100' },
  { code: '27', label: 'Combustíveis minerais, óleos minerais e produtos da destilação', defaultNcm8d: '27090010' },

  // Seção VI — Produtos das Indústrias Químicas
  { code: '28', label: 'Produtos químicos inorgânicos', defaultNcm8d: '28369200' },
  { code: '29', label: 'Produtos químicos orgânicos', defaultNcm8d: '29053100' },
  { code: '30', label: 'Produtos farmacêuticos', defaultNcm8d: '30049099' },
  { code: '31', label: 'Adubos e fertilizantes', defaultNcm8d: '31021000' },
  { code: '32', label: 'Extratos tanantes e tintoriais; tintas e vernizes', defaultNcm8d: '32089090' },
  { code: '33', label: 'Óleos essenciais; produtos de perfumaria e cosméticos', defaultNcm8d: '33012900' },
  { code: '34', label: 'Sabões, detergentes e preparações para lavagem', defaultNcm8d: '34021190' },
  { code: '35', label: 'Matérias albuminoides; amidos modificados; colas; enzimas', defaultNcm8d: '35011000' },
  { code: '36', label: 'Pólvoras e explosivos; artigos de pirotecnia', defaultNcm8d: '36010000' },
  { code: '37', label: 'Produtos para uso fotográfico ou cinematográfico', defaultNcm8d: '37019190' },
  { code: '38', label: 'Produtos diversos das indústrias químicas', defaultNcm8d: '38089419' },

  // Seção VII — Plásticos e Borracha
  { code: '39', label: 'Plásticos e suas obras', defaultNcm8d: '39011090' },
  { code: '40', label: 'Borracha e suas obras', defaultNcm8d: '40011000' },

  // Seção VIII — Peles, Couros e Peleteria
  { code: '41', label: 'Peles e couros', defaultNcm8d: '41041911' },
  { code: '42', label: 'Obras de couro; artigos de seleiro e correeiro', defaultNcm8d: '42021210' },
  { code: '43', label: 'Peles com pelo e suas obras; peles artificiais', defaultNcm8d: '43021910' },

  // Seção IX — Madeira e Obras de Madeira
  { code: '44', label: 'Madeira, carvão vegetal e obras de madeira', defaultNcm8d: '44071190' },
  { code: '45', label: 'Cortiça e suas obras', defaultNcm8d: '45019000' },
  { code: '46', label: 'Obras de espartaria ou de cestaria', defaultNcm8d: '46021990' },

  // Seção X — Pastas de Madeira; Papel e Cartão
  { code: '47', label: 'Pastas de madeira e outras matérias fibrosas celulósicas', defaultNcm8d: '47032100' },
  { code: '48', label: 'Papel, cartão e suas obras', defaultNcm8d: '48101900' },
  { code: '49', label: 'Produtos editoriais, imprensa e artes gráficas', defaultNcm8d: '49019900' },

  // Seção XI — Matérias Têxteis e suas Obras
  { code: '50', label: 'Seda', defaultNcm8d: '50020000' },
  { code: '51', label: 'Lã, pelos finos ou grosseiros; fios e tecidos de crina', defaultNcm8d: '51011910' },
  { code: '52', label: 'Algodão', defaultNcm8d: '52010000' },
  { code: '53', label: 'Outras fibras têxteis vegetais; fios e tecidos de papel', defaultNcm8d: '53010010' },
  { code: '54', label: 'Filamentos sintéticos ou artificiais', defaultNcm8d: '54021990' },
  { code: '55', label: 'Fibras sintéticas ou artificiais descontínuas', defaultNcm8d: '55020000' },
  { code: '56', label: 'Pastas, feltros, não-tecidos; fios especiais; cordéis', defaultNcm8d: '56021090' },
  { code: '57', label: 'Tapetes e revestimentos para pisos de matérias têxteis', defaultNcm8d: '57023200' },
  { code: '58', label: 'Tecidos especiais; superfícies têxteis tufadas', defaultNcm8d: '58042000' },
  { code: '59', label: 'Tecidos impregnados, revestidos ou estratificados', defaultNcm8d: '59119090' },
  { code: '60', label: 'Tecidos de malha', defaultNcm8d: '60019200' },
  { code: '61', label: 'Vestuário e acessórios de malha', defaultNcm8d: '61091000' },
  { code: '62', label: 'Vestuário e acessórios, exceto de malha', defaultNcm8d: '62034200' },
  { code: '63', label: 'Outros artefatos têxteis confeccionados', defaultNcm8d: '63025100' },

  // Seção XII — Calçados, Chapéus e Artefatos de Uso Semelhante
  { code: '64', label: 'Calçados, polainas e artefatos semelhantes', defaultNcm8d: '64041100' },
  { code: '65', label: 'Chapéus e artefatos de uso semelhante', defaultNcm8d: '65061090' },
  { code: '66', label: 'Guarda-chuvas, sombrinhas e bengalas', defaultNcm8d: '66011000' },
  { code: '67', label: 'Penas preparadas; flores artificiais; obras de cabelo humano', defaultNcm8d: '67021000' },

  // Seção XIII — Obras de Pedra, Gesso, Cerâmica e Vidro
  { code: '68', label: 'Obras de pedra, gesso, cimento, amianto, mica e análogas', defaultNcm8d: '68022100' },
  { code: '69', label: 'Produtos cerâmicos', defaultNcm8d: '69072100' },
  { code: '70', label: 'Vidro e suas obras', defaultNcm8d: '70193200' },

  // Seção XIV — Pérolas, Pedras Preciosas e Metais Preciosos
  { code: '71', label: 'Pérolas, pedras preciosas, metais preciosos e suas obras', defaultNcm8d: '71081200' },

  // Seção XV — Metais Comuns e suas Obras
  { code: '72', label: 'Ferro fundido, ferro e aço', defaultNcm8d: '72082700' },
  { code: '73', label: 'Obras de ferro fundido, ferro ou aço', defaultNcm8d: '73042310' },
  { code: '74', label: 'Cobre e suas obras', defaultNcm8d: '74031100' },
  { code: '75', label: 'Níquel e suas obras', defaultNcm8d: '75030000' },
  { code: '76', label: 'Alumínio e suas obras', defaultNcm8d: '76011000' },
  { code: '78', label: 'Chumbo e suas obras', defaultNcm8d: '78011000' },
  { code: '79', label: 'Zinco e suas obras', defaultNcm8d: '79011100' },
  { code: '80', label: 'Estanho e suas obras', defaultNcm8d: '80011000' },
  { code: '81', label: 'Outros metais comuns; ceramais (cermets) e suas obras', defaultNcm8d: '81089090' },
  { code: '82', label: 'Ferramentas e instrumentos; artigos de cutelaria', defaultNcm8d: '82079010' },
  { code: '83', label: 'Obras diversas de metais comuns', defaultNcm8d: '83040000' },

  // Seção XVI — Máquinas, Aparelhos e Equipamentos
  { code: '84', label: 'Máquinas, aparelhos e instrumentos mecânicos; reatores nucleares', defaultNcm8d: '84713012' },
  { code: '85', label: 'Máquinas, aparelhos e materiais elétricos; smartphones e TVs', defaultNcm8d: '85171231' },

  // Seção XVII — Material de Transporte
  { code: '86', label: 'Veículos e material para vias férreas', defaultNcm8d: '86090000' },
  { code: '87', label: 'Veículos automóveis, tratores e motocicletas', defaultNcm8d: '87032310' },
  { code: '88', label: 'Aeronaves e aparelhos espaciais', defaultNcm8d: '88024000' },
  { code: '89', label: 'Embarcações e estruturas flutuantes', defaultNcm8d: '89012090' },

  // Seção XVIII — Instrumentos de Óptica, Medida e Relojoaria
  { code: '90', label: 'Instrumentos de óptica, fotografia, medida e controle', defaultNcm8d: '90271000' },
  { code: '91', label: 'Aparelhos de relojoaria', defaultNcm8d: '91021200' },
  { code: '92', label: 'Instrumentos musicais', defaultNcm8d: '92011000' },

  // Seção XIX — Armas e Munições
  { code: '93', label: 'Armas e munições; suas partes e acessórios', defaultNcm8d: '93012000' },

  // Seção XX — Mercadorias e Produtos Diversos
  { code: '94', label: 'Móveis; artigos de cama, colchões e almofadas', defaultNcm8d: '94016900' },
  { code: '95', label: 'Brinquedos, jogos e artigos para esportes', defaultNcm8d: '95030090' },
  { code: '96', label: 'Obras diversas', defaultNcm8d: '96081000' },

  // Seção XXI — Objetos de Arte e Antiguidades
  { code: '97', label: 'Objetos de arte, de coleção e antiguidades', defaultNcm8d: '97010000' },
];

/** Busca um capítulo NCM pelo código de 2 dígitos. Retorna undefined se não encontrar. */
export function getNcmChapter(code: string): NcmChapter | undefined {
  return NCM_CHAPTERS.find((c) => c.code === code);
}
