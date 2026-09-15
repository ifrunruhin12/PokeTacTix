import { motion } from 'framer-motion';
import PropTypes from 'prop-types';

/**
 * ItemCard Component
 * Displays a purchasable game item (e.g. evolution stone, booster) in the shop
 */
function ItemCard({ item, onPurchase, onUse, usingItem, userCoins, ownedQuantity }) {
  const canAfford = userCoins >= item.price;
  const isEvolutionItem = item.item_type === 'evolution';
  const isBooster = item.item_type === 'booster';
  const canUse = isBooster && ownedQuantity > 0;

  const getBorderClass = () => {
    if (isEvolutionItem) {
      return 'border-amber-400 shadow-amber-400/30';
    }
    if (isBooster) {
      return 'border-blue-400 shadow-blue-400/30';
    }
    return 'border-gray-600 shadow-gray-600/30';
  };

  return (
    <motion.div
      className={`relative bg-gray-800 rounded-lg border-2 ${getBorderClass()} overflow-hidden shadow-lg flex flex-col`}
      whileHover={{ scale: 1.05, y: -5 }}
      transition={{ duration: 0.2 }}
    >
      {/* Category badge */}
      <div className="absolute top-2 left-2 z-10">
        <span className={`${isEvolutionItem ? 'bg-amber-500' : 'bg-blue-500'} text-white text-xs font-bold px-2 py-1 rounded-full`}>
          {isEvolutionItem ? '⚡ EVOLUTION' : 'BOOSTER'}
        </span>
      </div>

      {/* Owned quantity */}
      {ownedQuantity > 0 && (
        <div className="absolute top-2 right-2 z-10">
          <span className="bg-green-600 text-white text-xs font-bold px-2 py-1 rounded-full">
            ✓ x{ownedQuantity}
          </span>
        </div>
      )}

      {/* Item icon */}
      <div className="h-32 flex items-center justify-center bg-gray-700/50">
        {item.icon ? (
          <img
            src={item.icon}
            alt={item.name}
            className="h-20 w-20 object-contain"
            loading="lazy"
            onError={(e) => { e.target.style.display = 'none'; }}
          />
        ) : (
          <span className="text-4xl">💎</span>
        )}
      </div>

      {/* Item details */}
      <div className="p-3 flex flex-col flex-1">
        <h3 className="text-white font-bold text-sm mb-1 truncate" title={item.name}>
          {item.name}
        </h3>
        <p className="text-gray-400 text-xs mb-3 line-clamp-2 flex-1" title={item.description}>
          {item.description}
        </p>

        <div className="flex items-center justify-between mb-2">
          <span className="text-yellow-400 font-bold">
            {item.price} 🪙
          </span>
        </div>

        <button
          onClick={() => onPurchase(item)}
          disabled={!canAfford}
          className={`w-full py-2 px-3 rounded-lg font-bold text-sm transition-colors ${
            canAfford
              ? 'bg-blue-600 hover:bg-blue-700 text-white'
              : 'bg-gray-600 text-gray-400 cursor-not-allowed'
          }`}
        >
          {!canAfford ? 'Not enough coins' : 'Purchase'}
        </button>

        {isBooster && (
          <button
            onClick={() => onUse(item)}
            disabled={!canUse || usingItem === item.id}
            className={`w-full mt-2 py-2 px-3 rounded-lg font-bold text-sm transition-colors ${
              canUse && usingItem !== item.id
                ? 'bg-green-600 hover:bg-green-700 text-white'
                : 'bg-gray-600 text-gray-400 cursor-not-allowed'
            }`}
            title={canUse ? 'Activate the booster for your next battles' : 'Buy the booster first to use it'}
          >
            {usingItem === item.id ? 'Activating...' : canUse ? 'Use' : 'Need to own it'}
          </button>
        )}
      </div>
    </motion.div>
  );
}

ItemCard.propTypes = {
  item: PropTypes.shape({
    id: PropTypes.string.isRequired,
    name: PropTypes.string.isRequired,
    description: PropTypes.string,
    item_type: PropTypes.string.isRequired,
    price: PropTypes.number.isRequired,
    icon: PropTypes.string,
  }).isRequired,
  onPurchase: PropTypes.func.isRequired,
  onUse: PropTypes.func,
  usingItem: PropTypes.string,
  userCoins: PropTypes.number.isRequired,
  ownedQuantity: PropTypes.number,
};

/**
 * ActiveBoostsPanel Component
 * Shows the player's currently active deck buffs and remaining battles
 */
function ActiveBoostsPanel({ boosts }) {
  if (!boosts || boosts.length === 0) return null;

  return (
    <div className="mb-6 bg-gray-800 rounded-lg border border-blue-500/40 p-4">
      <h3 className="text-white font-bold mb-3">💪 Active Boosts</h3>
      <div className="flex flex-wrap gap-2">
        {boosts.map(boost => (
          <div
            key={boost.id}
            className="bg-blue-500/20 border border-blue-500 rounded-lg px-3 py-2 flex items-center gap-2"
          >
            <span className="text-lg">{boost.stat === 'hp' ? '❤️' : '⚔️'}</span>
            <div>
              <p className="text-white text-sm font-bold">+{boost.bonus} {boost.stat.toUpperCase()}</p>
              <p className="text-gray-400 text-xs">{boost.item_name} · {boost.battles_remaining} battle{boost.battles_remaining !== 1 ? 's' : ''} left</p>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}

ActiveBoostsPanel.propTypes = {
  boosts: PropTypes.arrayOf(PropTypes.shape({
    id: PropTypes.number.isRequired,
    item_id: PropTypes.string,
    item_name: PropTypes.string,
    stat: PropTypes.string,
    bonus: PropTypes.number,
    battles_remaining: PropTypes.number,
  })),
};

/**
 * ItemsGrid Component
 * Displays the shop's game item category (evolution stones, boosters)
 */
export default function ItemsGrid({ items, onPurchase, onUse, usingItem, userCoins, inventory, activeBoosts }) {
  if (!items || items.length === 0) {
    return (
      <div className="text-center py-12 bg-gray-800 rounded-lg">
        <p className="text-gray-400 text-lg">No items available in the shop right now</p>
      </div>
    );
  }

  const ownedQuantity = (itemId) =>
    inventory?.find(entry => entry.id === itemId)?.quantity || 0;

  return (
    <div>
      <ActiveBoostsPanel boosts={activeBoosts} />
      <h2 className="text-2xl font-bold mb-4">
        Game Items ({items.length})
      </h2>
      <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 gap-4">
        {items.map(item => (
          <ItemCard
            key={item.id}
            item={item}
            onPurchase={onPurchase}
            onUse={onUse}
            usingItem={usingItem}
            userCoins={userCoins}
            ownedQuantity={ownedQuantity(item.id)}
          />
        ))}
      </div>
    </div>
  );
}

ItemsGrid.propTypes = {
  items: PropTypes.arrayOf(PropTypes.object).isRequired,
  onPurchase: PropTypes.func.isRequired,
  onUse: PropTypes.func,
  usingItem: PropTypes.string,
  userCoins: PropTypes.number.isRequired,
  inventory: PropTypes.arrayOf(PropTypes.object),
  activeBoosts: PropTypes.arrayOf(PropTypes.object),
};
